// Package infrastructure stores all implementations of app interfaces
package infrastructure

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
	"net"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/mappers"
	"rol/app/utils"
	"rol/domain"
	"strings"
)

//HostNetworkManager is a struct for network manager
type HostNetworkManager struct {
	configStorage     interfaces.IHostNetworkConfigStorage
	iptables          *iptables.IPTables
	hasUnsavedChanges bool
}

var netfilterTables = []string{
	"filter",
	"nat",
	"mangle",
	"raw",
	"security",
}

//NewHostNetworkManager constructor for HostNetworkManager
func NewHostNetworkManager(configStorage interfaces.IHostNetworkConfigStorage) (interfaces.IHostNetworkManager, error) {
	ipTables, err := iptables.New()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting iptables instance")
	}
	hostNetworkManager := &HostNetworkManager{
		configStorage:     configStorage,
		iptables:          ipTables,
		hasUnsavedChanges: true,
		// we set this flag for calling reset changes function at start, for apply configuration from storage
	}
	//if it's a first time, we need to save config based on current configuration
	_, err = configStorage.GetConfig()
	if err != nil && !errors.As(err, errors.Internal) {
		err = hostNetworkManager.SaveConfiguration()
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to save initial host network configuration to storage")
		}
	} else if errors.As(err, errors.Internal) {
		return nil, errors.Wrap(err, "failed to get configuration from storage")
	}
	//load configuration from storage
	err = hostNetworkManager.ResetChanges()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to load initial host network configuration from storage")
	}
	return hostNetworkManager, nil
}

func (h *HostNetworkManager) parseLinkAddr(link netlink.Link) ([]net.IPNet, error) {
	addrList, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "get addresses error")
	}

	var out []net.IPNet
	for _, addr := range addrList {
		out = append(out, *addr.IPNet)
	}
	return out, nil
}

func (h *HostNetworkManager) getParentName(link netlink.Link) (string, error) {
	parent, err := netlink.LinkByIndex(link.Attrs().ParentIndex)
	if err != nil {
		return "", errors.Internal.Wrap(err, "get host interface by index failed")
	}
	return parent.Attrs().Name, nil
}

func (h *HostNetworkManager) getSlaves(master netlink.Link) ([]string, error) {
	out := []string{}
	links, err := netlink.LinkList()
	if err != nil {
		return out, errors.Internal.Wrap(err, "error getting a list of link devices")
	}
	for _, link := range links {
		if link.Attrs().MasterIndex == master.Attrs().Index {
			out = append(out, link.Attrs().Name)
		}
	}
	return out, nil
}

func (h *HostNetworkManager) mapLink(link netlink.Link) (interfaces.IHostNetworkLink, error) {
	if link.Type() == "device" {
		addresses, err := h.parseLinkAddr(link)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "error parsing link addresses")
		}
		device := domain.HostNetworkDevice{
			HostNetworkLink: domain.HostNetworkLink{
				Name:      link.Attrs().Name,
				Type:      link.Type(),
				Addresses: addresses,
			}}
		return device, nil

	} else if link.Type() == "vlan" {
		addresses, err := h.parseLinkAddr(link)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "error parsing link addresses")
		}
		parent, err := h.getParentName(link)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "error getting parent name")
		}
		vlan := domain.HostNetworkVlan{
			HostNetworkLink: domain.HostNetworkLink{
				Name:      link.Attrs().Name,
				Type:      link.Type(),
				Addresses: addresses,
			},
			VlanID: link.(*netlink.Vlan).VlanId,
			Parent: parent,
		}
		return vlan, nil
	} else if link.Type() == "bridge" {
		addresses, err := h.parseLinkAddr(link)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "error parsing link addresses")
		}
		slaves, err := h.getSlaves(link)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "get slaves failed")
		}
		bridge := domain.HostNetworkBridge{
			HostNetworkLink: domain.HostNetworkLink{
				Name:      link.Attrs().Name,
				Type:      link.Type(),
				Addresses: addresses,
			},
			Slaves: slaves,
		}
		return bridge, nil
	}
	return domain.HostNetworkLink{Name: link.Attrs().Name, Type: "none", Addresses: []net.IPNet{}}, nil
}

//GetList gets list of host network interfaces
//
//Return:
//	[]interfaces.IHostNetworkLink - list of interfaces
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) GetList() ([]interfaces.IHostNetworkLink, error) {
	var out []interfaces.IHostNetworkLink
	links, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error getting a list of link devices")
	}
	for _, link := range links {
		networkInterface, err := h.mapLink(link)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to map device link to HostNetworkLink")
		}
		out = append(out, networkInterface)
	}

	return out, nil
}

//GetByName gets host network interface by its name
//
//Params:
//	name - interface name
//Return:
//	interfaces.IHostNetworkLink - interfaces
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) GetByName(name string) (interfaces.IHostNetworkLink, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		if err.Error() == "Link not found" {
			return nil, errors.NotFound.New("link with this name is not exist")
		}
		return nil, errors.Internal.Wrap(err, "failed to get link by name")
	}
	if link == nil {
		return nil, errors.NotFound.New("link with this name is not exist")
	}
	out, err := h.mapLink(link)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to map device link to HostNetworkLink")
	}
	return out, nil
}

//CreateVlan creates vlan on host
//
//Params:
//	master - name of the master network interface
//	vlanID - ID to be set for vlan
//Return:
//	string - new vlan name that will be rol.{master}.{vlanID}
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) CreateVlan(master string, vlanID int) (string, error) {
	parent, err := netlink.LinkByName(master)
	if err != nil {
		return "", errors.Internal.Wrap(err, "getting device link by name failed")
	}
	la := netlink.NewLinkAttrs()

	vlanName := fmt.Sprintf("rol.%s.%d", master, vlanID)
	la.Name = vlanName
	la.ParentIndex = parent.Attrs().Index
	vlan := &netlink.Vlan{
		LinkAttrs:    la,
		VlanId:       vlanID,
		VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
	}
	err = netlink.LinkAdd(vlan)
	if err != nil {
		return "", errors.Internal.Wrap(err, "failed to add vlan link")
	}

	h.hasUnsavedChanges = true

	return vlanName, nil
}

//CreateBridge creates bridge on host
//
//Params:
//	name - new bridge name
//Return:
//	string - new bridge name that will be rol.br.{name}
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) CreateBridge(name string) (string, error) {
	la := netlink.NewLinkAttrs()
	bridgeName := fmt.Sprintf("rol.br.%s", name)
	la.Name = bridgeName
	bridge := &netlink.Bridge{LinkAttrs: la}
	err := netlink.LinkAdd(bridge)
	if err != nil {
		return "", errors.Internal.Wrap(err, "failed to add bridge link")
	}
	h.hasUnsavedChanges = true
	return bridgeName, nil
}

//SetLinkMaster set master for link
//
//Params:
//	slaveName - name of link that will be slave
//	masterName - name of link that will be master for the slave
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) SetLinkMaster(slaveName, masterName string) error {
	slave, err := netlink.LinkByName(slaveName)
	if err != nil {
		return errors.Internal.Wrap(err, "getting slave by name failed")
	}
	master, err := netlink.LinkByName(masterName)
	if err != nil {
		return errors.Internal.Wrap(err, "getting master by name failed")
	}
	err = netlink.LinkSetMaster(slave, master)
	if err != nil {
		return errors.Internal.Wrap(err, "failed set link master")
	}
	h.hasUnsavedChanges = true
	return nil
}

//UnsetLinkMaster removes the master of the link
//
//Params:
//	linkName - name of the link
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) UnsetLinkMaster(linkName string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return errors.Internal.Wrap(err, "getting link by name failed")
	}
	err = netlink.LinkSetNoMaster(link)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to set no master for link")
	}
	h.hasUnsavedChanges = true
	return nil
}

//SetLinkUp enables the link
//
//Params:
//	linkName - name of the link
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) SetLinkUp(linkName string) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return errors.Internal.Wrap(err, "getting link by name failed")
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		return errors.Internal.Wrap(err, "link set up failed")
	}
	return nil
}

//DeleteLinkByName deletes interface on host by its name
//
//Params:
//	name - interface name
//Return
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) DeleteLinkByName(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return errors.Internal.Wrap(err, "getting link by name failed")
	}
	err = netlink.LinkDel(link)
	if err != nil {
		return errors.Internal.Wrap(err, "deleting link failed")
	}
	h.hasUnsavedChanges = true
	return nil
}

//AddrAdd Add new ip address for network interface
//
//Params:
//	linkName - name of the interface
//	addr - ip address with mask net.IPNet
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) AddrAdd(linkName string, addr net.IPNet) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return errors.Internal.Wrap(err, "getting link by name failed")
	}
	cidr := addr.String()
	linkAddr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return errors.Internal.Wrap(err, "parse cidr address failed")
	}
	err = netlink.AddrAdd(link, linkAddr)
	if err != nil {
		return errors.Internal.Wrap(err, "error adding address to link")
	}
	h.hasUnsavedChanges = true
	return nil
}

//AddrDelete Delete ip address from network interface
//
//Params:
//	linkName - name of the interface
//	addr - ip address with mask net.IPNet
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) AddrDelete(linkName string, addr net.IPNet) error {
	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return errors.Internal.Wrap(err, "getting link by name failed")
	}
	cidr := addr.String()
	linkAddr, err := netlink.ParseAddr(cidr)
	if err != nil {
		return errors.Internal.Wrap(err, "parse cidr address failed")
	}
	err = netlink.AddrDel(link, linkAddr)
	if err != nil {
		return errors.Internal.Wrap(err, "error adding address to link")
	}
	h.hasUnsavedChanges = true
	return nil
}

func (h *HostNetworkManager) parseRule(rule domain.HostNetworkTrafficRule) []string {
	var rulespec []string
	if rule.Source != "" {
		rulespec = append(rulespec, "-s")
		rulespec = append(rulespec, rule.Source)
	}
	if rule.Destination != "" {
		rulespec = append(rulespec, "-d")
		rulespec = append(rulespec, rule.Destination)
	}
	if rule.Action != "" {
		rulespec = append(rulespec, "-j")
		rulespec = append(rulespec, rule.Action)
	}
	return rulespec
}

//CreateTrafficRule Create netfilter traffic rule for specified table
//
//Params:
//	table - table to create a rule
//	rule - rule entity
//Return:
//	domain.HostNetworkTrafficRule - new traffic rule
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) CreateTrafficRule(table string, rule domain.HostNetworkTrafficRule) (domain.HostNetworkTrafficRule, error) {
	rulespec := h.parseRule(rule)
	err := h.iptables.AppendUnique(table, rule.Chain, rulespec...)
	if err != nil {
		return domain.HostNetworkTrafficRule{}, errors.Internal.Wrap(err, "failed to create traffic rule")
	}
	h.hasUnsavedChanges = true
	return rule, nil
}

//DeleteTrafficRule Delete netfilter traffic rule in specified table
//
//Params:
//	table - table to delete a rule
//	rule - rule entity
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) DeleteTrafficRule(table string, rule domain.HostNetworkTrafficRule) error {
	rulespec := h.parseRule(rule)
	exist, err := h.iptables.Exists(table, rule.Chain, rulespec...)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to check existence of traffic rule")
	}
	if exist {
		err = h.iptables.Delete(table, rule.Chain, rulespec...)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to delete traffic rule")
		}
		h.hasUnsavedChanges = true
		return nil
	}
	return errors.NotFound.New("traffic rule not found")
}

//GetChainRules Get selected netfilter chain rules at specified table
//
//Params:
//	table - table to get a rules
//	chain - chain where we get the rules
//Return:
//	[]domain.HostNetworkTrafficRule - slice of rules
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) GetChainRules(table string, chain string) ([]domain.HostNetworkTrafficRule, error) {
	var rules []domain.HostNetworkTrafficRule

	list, err := h.iptables.Stats(table, chain)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to get list of traffic rules")
	}
	for _, l := range list {
		stat, err := h.iptables.ParseStat(l)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to parse traffic rule to stat struct")
		}
		rule := &domain.HostNetworkTrafficRule{Chain: chain}
		mappers.MapStatToTrafficRule(stat, rule)
		rules = append(rules, *rule)
	}
	return rules, err
}

func (h *HostNetworkManager) trimExclamationMarkInStat(slice []string) (out []string) {
	for _, element := range slice {
		out = append(out, strings.Trim(element, "!"))
	}
	return
}

//GetTableRules Get specified netfilter table rules
//
//Params:
//	table - table to get a rules
//Return:
//	[]domain.HostNetworkTrafficRule - slice of rules
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) GetTableRules(table string) ([]domain.HostNetworkTrafficRule, error) {
	var rules []domain.HostNetworkTrafficRule

	chains, err := h.iptables.ListChains(table)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to get list table chains")
	}
	for _, chain := range chains {
		list, err := h.iptables.Stats(table, chain)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to get list of traffic rules")
		}
		for _, l := range list {
			stat, err := h.iptables.ParseStat(h.trimExclamationMarkInStat(l))
			if err != nil {
				return nil, errors.Internal.Wrap(err, "failed to parse traffic rule to stat struct")
			}
			rule := &domain.HostNetworkTrafficRule{Chain: chain}
			mappers.MapStatToTrafficRule(stat, rule)
			rules = append(rules, *rule)
		}
	}
	return rules, nil
}

//SaveConfiguration save current host network configuration to the configuration storage
//Save previous config file to .back file
//
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) SaveConfiguration() error {
	config := domain.HostNetworkConfig{}
	networkInterfaces, err := h.GetList()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get list of host network interfaces")
	}
	for _, inter := range networkInterfaces {
		if inter.GetType() == "vlan" {
			config.Vlans = append(config.Vlans, inter.(domain.HostNetworkVlan))
		} else if inter.GetType() == "device" {
			config.Devices = append(config.Devices, inter.(domain.HostNetworkDevice))
		} else if inter.GetType() == "bridge" {
			config.Bridges = append(config.Bridges, inter.(domain.HostNetworkBridge))
		}
	}
	for _, table := range netfilterTables {
		rules, err := h.GetTableRules(table)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to get table rules")
		}

		h.setTrafficRulesConfigField(table, rules, &config)
	}
	err = h.configStorage.SaveConfig(config)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to save host network config to storage")
	}
	h.hasUnsavedChanges = false
	return nil
}

func (h *HostNetworkManager) bridgeExistOnHost(links []interfaces.IHostNetworkLink, bridgeName string) bool {
	for _, inter := range links {
		if inter.GetType() == "bridge" && inter.GetName() == bridgeName {
			return true
		}
	}
	return false
}

func (h *HostNetworkManager) bridgeExistInConfig(config domain.HostNetworkConfig, bridgeName string) bool {
	for _, bridge := range config.Bridges {
		if bridge.GetName() == bridgeName {
			return true
		}
	}
	return false
}

func (h *HostNetworkManager) vlanExistOnHost(links []interfaces.IHostNetworkLink, vlanName string) bool {
	for _, inter := range links {
		if inter.GetType() == "vlan" && inter.GetName() == vlanName {
			return true
		}
	}
	return false
}

func (h *HostNetworkManager) vlanExistInConfig(config domain.HostNetworkConfig, vlanName string) bool {
	for _, vlan := range config.Vlans {
		if vlan.GetName() == vlanName {
			return true
		}
	}
	return false
}

func (h *HostNetworkManager) addressExistOnHostLink(links []interfaces.IHostNetworkLink, linkName string, address net.IPNet) bool {
	for _, inter := range links {
		if inter.GetName() == linkName {
			addresses := inter.GetAddresses()
			for _, addr := range addresses {
				if addr.String() == address.String() {
					return true
				}
			}
		}
	}
	return false
}

func (h *HostNetworkManager) addressExistInLinkConfig(config domain.HostNetworkConfig, linkName string, address net.IPNet) bool {
	for _, inter := range config.Vlans {
		if inter.GetName() == linkName {
			addresses := inter.GetAddresses()
			for _, addr := range addresses {
				if addr.String() == address.String() {
					return true
				}
			}
		}
	}
	for _, inter := range config.Devices {
		if inter.GetName() == linkName {
			addresses := inter.GetAddresses()
			for _, addr := range addresses {
				if addr.String() == address.String() {
					return true
				}
			}
		}
	}
	for _, inter := range config.Bridges {
		if inter.GetName() == linkName {
			addresses := inter.GetAddresses()
			for _, addr := range addresses {
				if addr.String() == address.String() {
					return true
				}
			}
		}
	}
	//TODO: If we add new type of the interfaces, we must not forget to add it here.
	return false
}

func (h *HostNetworkManager) loadVlanConfiguration(config domain.HostNetworkConfig) error {
	hostLinks, err := h.GetList()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get list of host network interfaces")
	}
	// Add all settings from config to vlan interfaces
	for _, vlan := range config.Vlans {
		//Skip all vlans that not configured by our system
		if !strings.Contains(vlan.Name, "rol.") {
			continue
		}
		vlanExist := h.vlanExistOnHost(hostLinks, vlan.Name)
		if !vlanExist {
			vlanName, err := h.CreateVlan(vlan.Parent, vlan.VlanID)
			if err != nil {
				return errors.Internal.Wrap(err, "error when creating a vlan")
			}
			for _, addr := range vlan.Addresses {
				err = h.AddrAdd(vlanName, addr)
				if err != nil {
					return errors.Internal.Wrap(err, "failed set address to vlan")
				}
			}
		} else {
			for _, addr := range vlan.Addresses {
				h.addressExistOnHostLink(hostLinks, vlan.GetName(), addr)
				err = h.AddrAdd(vlan.GetName(), addr)
				if err != nil {
					return errors.Internal.Wrap(err, "failed set address to vlan")
				}
			}
		}
	}
	// Remove all configurations that not exist in config
	networkInterfaces, err := h.GetList()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get list of host network interfaces")
	}
	for _, inter := range networkInterfaces {
		if !strings.Contains(inter.GetName(), "rol.") && inter.GetType() != "vlan" {
			continue
		}
		if !h.vlanExistInConfig(config, inter.GetName()) {
			err := h.DeleteLinkByName(inter.GetName())
			if err != nil {
				return errors.Internal.Wrap(err, "delete link by name error")
			}
		} else if h.vlanExistInConfig(config, inter.GetName()) {
			addresses := inter.GetAddresses()
			for _, address := range addresses {
				if !h.addressExistInLinkConfig(config, inter.GetName(), address) {
					err = h.AddrDelete(inter.GetName(), address)
					if err != nil {
						return errors.Internal.Wrap(err, "failed delete address from vlan")
					}
				}
			}
		}
	}
	return nil
}

func (h *HostNetworkManager) loadBridgeConfiguration(config domain.HostNetworkConfig) error {
	hostLinks, err := h.GetList()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to get list of host network interfaces")
	}
	// Add all settings from config to bridge interfaces
	for _, bridge := range config.Bridges {
		//Skip all bridges that not configured by our system
		if !strings.Contains(bridge.Name, "rol.br.") {
			continue
		}
		bridgeExist := h.bridgeExistOnHost(hostLinks, bridge.Name)
		if !bridgeExist {
			bridgeName, err := h.CreateBridge(bridge.Name[7:])
			if err != nil {
				return errors.Internal.Wrap(err, "error when creating a bridge")
			}
			for _, addr := range bridge.Addresses {
				err = h.AddrAdd(bridgeName, addr)
				if err != nil {
					return errors.Internal.Wrap(err, "failed set address to bridge")
				}
			}
		} else {
			for _, addr := range bridge.Addresses {
				h.addressExistOnHostLink(hostLinks, bridge.GetName(), addr)
				err = h.AddrAdd(bridge.GetName(), addr)
				if err != nil {
					return errors.Internal.Wrap(err, "failed set address to bridge")
				}
			}
		}
	}
	// Remove all configurations that not exist in config
	for _, inter := range hostLinks {
		if !strings.Contains(inter.GetName(), "rol.br.") && inter.GetType() != "bridge" {
			continue
		}
		if !h.bridgeExistInConfig(config, inter.GetName()) {
			err := h.DeleteLinkByName(inter.GetName())
			if err != nil {
				return errors.Internal.Wrap(err, "delete link by name error")
			}
		} else if h.bridgeExistInConfig(config, inter.GetName()) {
			addresses := inter.GetAddresses()
			for _, address := range addresses {
				if !h.addressExistInLinkConfig(config, inter.GetName(), address) {
					err = h.AddrDelete(inter.GetName(), address)
					if err != nil {
						return errors.Internal.Wrap(err, "failed delete address from vlan")
					}
				}
			}
		}
	}
	return nil
}

func (h *HostNetworkManager) setTrafficRulesConfigField(table string, rule []domain.HostNetworkTrafficRule, config *domain.HostNetworkConfig) {
	switch table {
	case "filter":
		config.TrafficRules.Filter = rule
	case "nat":
		config.TrafficRules.NAT = rule
	case "mangle":
		config.TrafficRules.Mangle = rule
	case "raw":
		config.TrafficRules.Raw = rule
	case "security":
		config.TrafficRules.Security = rule
	}
}

func (h *HostNetworkManager) getTrafficRulesConfigField(table string, config domain.HostNetworkConfig) []domain.HostNetworkTrafficRule {
	switch table {
	case "filter":
		return config.TrafficRules.Filter
	case "nat":
		return config.TrafficRules.NAT
	case "mangle":
		return config.TrafficRules.Mangle
	case "raw":
		return config.TrafficRules.Raw
	case "security":
		return config.TrafficRules.Security
	default:
		return nil
	}
}

func (h *HostNetworkManager) loadTrafficConfiguration(config domain.HostNetworkConfig) error {
	for _, table := range netfilterTables {
		rules, err := h.GetTableRules(table)
		if err != nil {
			return errors.Internal.Wrap(err, "failed to get table rules")
		}

		configField := h.getTrafficRulesConfigField(table, config)
		if configField == nil {
			return errors.Internal.New("failed to get config field")
		}

		for _, rule := range configField {
			if !utils.SliceContainsElement(rules, rule) {
				_, err = h.CreateTrafficRule(table, rule)
				if err != nil {
					return errors.Internal.New("error when creating traffic rule")
				}
			}
		}
		for _, rule := range rules {
			if !utils.SliceContainsElement(configField, rule) {
				err = h.DeleteTrafficRule(table, rule)
				if err != nil {
					return errors.Internal.New("error when deleting traffic rule")
				}
			}
		}
	}
	return nil
}

func (h *HostNetworkManager) loadConfiguration(config domain.HostNetworkConfig) error {
	err := h.loadVlanConfiguration(config)
	if err != nil {
		return errors.Internal.Wrap(err, "error loading vlan configuration")
	}
	err = h.loadBridgeConfiguration(config)
	if err != nil {
		return errors.Internal.Wrap(err, "error loading bridge configuration")
	}
	err = h.loadTrafficConfiguration(config)
	if err != nil {
		return errors.Internal.Wrap(err, "error loading host network traffic configuration")
	}
	h.hasUnsavedChanges = false
	return nil
}

//ResetChanges Reset all applied changes to state from saved configuration
//
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) ResetChanges() error {
	if h.hasUnsavedChanges == true {
		config, err := h.configStorage.GetConfig()
		if err != nil {
			return errors.Internal.Wrap(err, "error while getting network configuration from storage")
		}
		err = h.loadConfiguration(config)
		if err != nil {
			return errors.Internal.Wrap(err, "load backup configuration failed")
		}
		h.hasUnsavedChanges = false
	}
	return nil
}

//RestoreFromBackup restore and apply host network configuration from backup
//
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) RestoreFromBackup() error {
	backConfig, err := h.configStorage.GetBackupConfig()
	if err != nil {
		return errors.Internal.Wrap(err, "failed to restore host network configuration from backup")
	}
	err = h.loadConfiguration(backConfig)
	if err != nil {
		return errors.Internal.Wrap(err, "load backup configuration failed")
	}
	h.hasUnsavedChanges = true
	return nil
}

//HasUnsavedChanges Gets a flag about unsaved changes
//
//Return:
//	bool - if unsaved changes exist - true, otherwise false
func (h *HostNetworkManager) HasUnsavedChanges() bool {
	return h.hasUnsavedChanges
}
