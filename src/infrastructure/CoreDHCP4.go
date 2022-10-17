package infrastructure

import (
	"github.com/coredhcp/coredhcp/plugins"
	pluginDNS "github.com/insei/coredhcp/plugins/dns"
	pluginNetmask "github.com/insei/coredhcp/plugins/netmask"
	pluginRouter "github.com/insei/coredhcp/plugins/router"
	pluginServerid "github.com/insei/coredhcp/plugins/serverid"
	"net"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
	"strings"

	"github.com/coredhcp/coredhcp/config"
	"github.com/coredhcp/coredhcp/server"
)

var pluginsInitialized = false

func initializeCoreDHCPPlugins(leasesRepo interfaces.IGenericRepository[domain.DHCP4Lease]) error {
	pluginsSlice := []*plugins.Plugin{
		&pluginDNS.Plugin,
		&pluginNetmask.Plugin,
		NewRangeRepositoryPlugin(leasesRepo),
		&pluginRouter.Plugin,
		&pluginServerid.Plugin,
	}
	for _, plugin := range pluginsSlice {
		if err := plugins.RegisterPlugin(plugin); err != nil {
			return errors.Internal.Wrapf(err, "failed to register plugin: %v", plugin.Name)
		}
	}
	return nil
}

type coreDHCP4Server struct {
	config *config.Config
	server *server.Servers
	state  domain.DHCPServerState
}

//NewCoreDHCP4Server constructor for core DHCP v4 server
func NewCoreDHCP4Server(
	dhcp4config domain.DHCP4Config,
	leasesRepo interfaces.IGenericRepository[domain.DHCP4Lease],
) (interfaces.IDHCP4Server, error) {
	if !pluginsInitialized {
		err := initializeCoreDHCPPlugins(leasesRepo)
		if err != nil {
			return nil, err
		}
		pluginsInitialized = true
	}
	serv := &coreDHCP4Server{}
	err := serv.ReloadConfiguration(dhcp4config)
	if err != nil {
		return nil, err
	}
	return serv, nil
}

//ReloadConfiguration DHCP v4 server from config
func (s *coreDHCP4Server) ReloadConfiguration(dhcp4config domain.DHCP4Config) error {
	startEndIPs := strings.Split(dhcp4config.Range, "-")
	if len(startEndIPs) < 2 {
		return errors.Internal.Newf("incorrect ip range: %s", dhcp4config.Range)
	}
	s.config = &config.Config{
		Server6: nil,
		Server4: &config.ServerConfig{
			Addresses: []net.UDPAddr{{
				IP:   net.ParseIP("0.0.0.0"),
				Port: dhcp4config.Port,
				Zone: dhcp4config.Interface,
			}},
			Plugins: []config.PluginConfig{
				{
					Name: "dns",
					Args: strings.Split(dhcp4config.DNS, ";"),
				},
				{
					Name: "range_repo",
					Args: []string{
						dhcp4config.ID.String(),
						startEndIPs[0],
						startEndIPs[1],
						"3600s",
					},
				},
				{
					Name: "router",
					Args: []string{dhcp4config.Gateway},
				},
				{
					Name: "server_id",
					Args: []string{dhcp4config.ServerID},
				},
				{
					Name: "netmask",
					Args: []string{dhcp4config.Mask},
				},
			},
		},
	}
	return nil
}

//Start DHCP v4 server
func (s *coreDHCP4Server) Start() error {
	coredhcp, err := server.Start(s.config)
	if err != nil {
		s.state = domain.DHCPStateError
		return errors.Internal.Wrap(err, "failed to start dhcp v4 server")
	}
	s.state = domain.DHCPStateLaunched
	s.server = coredhcp
	return nil
}

//Stop DHCP v4 server
func (s *coreDHCP4Server) Stop() {
	if s.server != nil {
		s.server.Close()
	}
	s.state = domain.DHCPStateStopped
	s.server = nil
}

//GetState of DHCP v4 server
func (s *coreDHCP4Server) GetState() domain.DHCPServerState {
	return s.state
}
