package infrastructure

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"rol/app/interfaces"
	"rol/domain"
	"sync"
	"time"

	"github.com/coredhcp/coredhcp/handler"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/plugins"
	"github.com/coredhcp/coredhcp/plugins/allocators"
	"github.com/coredhcp/coredhcp/plugins/allocators/bitmap"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

var log = logger.GetLogger("plugins/range_repo")
var leasesRepo interfaces.IGenericRepository[uuid.UUID, domain.DHCP4Lease]

//NewRangeRepositoryPlugin constructor for range plugin that integrated with leases repository
func NewRangeRepositoryPlugin(repo interfaces.IGenericRepository[uuid.UUID, domain.DHCP4Lease]) *plugins.Plugin {
	leasesRepo = repo
	return &plugins.Plugin{
		Name:   "range_repo",
		Setup4: setupRange,
	}
}

//Record holds an IP lease record
type Record struct {
	ID      uuid.UUID
	IP      net.IP
	expires time.Time
}

// PluginState is the data held by an instance of the range plugin
type PluginState struct {
	// Rough lock for the whole plugin, we'll get better performance once we use leasestorage
	sync.Mutex
	LeaseTime time.Duration
	allocator allocators.Allocator
	serverID  uuid.UUID
}

func (p *PluginState) getLeaseFromRepo(mac string) (*Record, error) {
	if leasesRepo == nil {
		return nil, errors.New("repository is not set")
	}
	ctx := context.Background()
	queryBuilder := leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", p.serverID)
	queryBuilder.Where("MAC", "==", mac)
	leases, err := leasesRepo.GetList(ctx, "", "", 1, 1, queryBuilder)
	if err != nil {
		return nil, err
	}
	if len(leases) > 0 {
		return &Record{
			ID:      leases[0].ID,
			IP:      net.ParseIP(leases[0].IP),
			expires: leases[0].Expires,
		}, nil
	}
	return nil, nil
}

func (p *PluginState) createLeaseInRepo(addr net.HardwareAddr, rec *Record) error {
	newLease := domain.DHCP4Lease{
		IP:            rec.IP.String(),
		MAC:           addr.String(),
		Expires:       rec.expires,
		DHCP4ConfigID: p.serverID,
	}
	createdLease, err := leasesRepo.Insert(context.Background(), newLease)
	rec.ID = createdLease.ID
	return err
}

func (p *PluginState) updateLeaseExpiresTimeInRepo(rec *Record) error {
	ctx := context.Background()
	lease, err := leasesRepo.GetByID(ctx, rec.ID)
	if err != nil {
		return err
	}
	lease.Expires = rec.expires
	_, err = leasesRepo.Update(ctx, lease)
	if err != nil {
		return nil
	}
	return err
}

// Handler4 handles DHCPv4 packets for the range plugin
func (p *PluginState) Handler4(req, resp *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, bool) {
	p.Lock()
	defer p.Unlock()

	record, err := p.getLeaseFromRepo(req.ClientHWAddr.String())
	if err != nil {
		log.Errorf("failed to get ip address for mac %v from repository: %v", req.ClientHWAddr.String(), err)
		return nil, true
	}
	if record == nil {
		// Allocating new address since there isn't one allocated
		log.Printf("MAC address %s is new, leasing new IPv4 address", req.ClientHWAddr.String())
		ip, err := p.allocator.Allocate(net.IPNet{})
		if err != nil {
			log.Errorf("Could not allocate IP for MAC %s: %v", req.ClientHWAddr.String(), err)
			return nil, true
		}
		rec := Record{
			IP:      ip.IP.To4(),
			expires: time.Now().Add(p.LeaseTime),
		}
		err = p.createLeaseInRepo(req.ClientHWAddr, &rec)
		if err != nil {
			log.Errorf("SaveIPAddress for MAC %s failed: %v", req.ClientHWAddr.String(), err)
		}
		record = &rec
	} else {
		// Ensure we extend the existing lease at least past when the one we're giving expires
		if record.expires.Before(time.Now().Add(p.LeaseTime)) {
			record.expires = time.Now().Add(p.LeaseTime).Round(time.Second)
			err := p.updateLeaseExpiresTimeInRepo(record)
			if err != nil {
				log.Errorf("Could not persist lease for MAC %s: %v", req.ClientHWAddr.String(), err)
			}
		}
	}
	resp.YourIPAddr = record.IP
	resp.Options.Update(dhcpv4.OptIPAddressLeaseTime(p.LeaseTime.Round(time.Second)))
	log.Printf("found IP address %s for MAC %s", record.IP, req.ClientHWAddr.String())
	return resp, false
}

func loadRecordsFromRepo(serverID uuid.UUID) (map[string]*Record, error) {
	if leasesRepo == nil {
		return nil, errors.New("repository is not set")
	}
	ctx := context.Background()
	queryBuilder := leasesRepo.NewQueryBuilder(ctx)
	queryBuilder.Where("DHCP4ConfigID", "==", serverID)
	leasesCount, err := leasesRepo.Count(ctx, queryBuilder)
	if err != nil {
		return nil, err
	}
	leases, err := leasesRepo.GetList(ctx, "", "", 1, int(leasesCount), queryBuilder)
	if err != nil {
		return nil, err
	}
	records := make(map[string]*Record)
	for _, lease := range leases {
		records[lease.MAC] = &Record{
			ID:      lease.ID,
			IP:      net.ParseIP(lease.IP),
			expires: lease.Expires,
		}
	}
	return records, nil
}

func setupRange(args ...string) (handler.Handler4, error) {
	var (
		err error
		p   PluginState
	)

	if len(args) < 4 {
		return nil, fmt.Errorf("invalid number of arguments, want: 4 (file name, start IP, end IP, lease time), got: %d", len(args))
	}
	serverIDStr := args[0]
	if serverIDStr == "" {
		return nil, errors.New("server id cannot be empty")
	}
	serverID, err := uuid.Parse(serverIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid server ID: %v", serverIDStr)
	}
	p.serverID = serverID
	ipRangeStart := net.ParseIP(args[1])
	if ipRangeStart.To4() == nil {
		return nil, fmt.Errorf("invalid IPv4 address: %v", args[1])
	}
	ipRangeEnd := net.ParseIP(args[2])
	if ipRangeEnd.To4() == nil {
		return nil, fmt.Errorf("invalid IPv4 address: %v", args[2])
	}
	if binary.BigEndian.Uint32(ipRangeStart.To4()) >= binary.BigEndian.Uint32(ipRangeEnd.To4()) {
		return nil, errors.New("start of IP range has to be lower than the end of an IP range")
	}

	p.allocator, err = bitmap.NewIPv4Allocator(ipRangeStart, ipRangeEnd)
	if err != nil {
		return nil, fmt.Errorf("could not create an allocator: %w", err)
	}

	p.LeaseTime, err = time.ParseDuration(args[3])
	if err != nil {
		return nil, fmt.Errorf("invalid lease duration: %v", args[3])
	}

	recordsv4, err := loadRecordsFromRepo(p.serverID)
	if err != nil {
		return nil, fmt.Errorf("could not load records from repository: %v", err)
	}

	log.Printf("Loaded %d DHCPv4 leases from repository", len(recordsv4))

	for _, v := range recordsv4 {
		ip, err := p.allocator.Allocate(net.IPNet{IP: v.IP})
		if err != nil {
			return nil, fmt.Errorf("failed to re-allocate leased ip %v: %v", v.IP.String(), err)
		}
		if ip.IP.String() != v.IP.String() {
			return nil, fmt.Errorf("allocator did not re-allocate requested leased ip %v: %v", v.IP.String(), ip.String())
		}
	}
	return p.Handler4, nil
}
