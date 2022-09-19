//go:build linux

package tests

import (
	"net"
	"os"
	"path/filepath"
	"rol/app/interfaces"
	"rol/domain"
	"rol/infrastructure"
	"testing"
)

type networkManagerTester struct {
	bridgeName     string
	vlanName       string
	vlanID         int
	manager        interfaces.IHostNetworkManager
	storage        interfaces.IHostNetworkConfigStorage
	configFilePath string
}

var netManagerTester *networkManagerTester

func Test_HostNetworkManager_Prepare(t *testing.T) {
	filePath, _ := os.Executable()
	netManagerTester = &networkManagerTester{}

	netManagerTester.configFilePath = filepath.Join(filepath.Dir(filePath), "hostNetworkConfig.yaml")
	netManagerTester.storage = infrastructure.NewYamlHostNetworkConfigStorage(domain.GlobalDIParameters{
		RootPath: filepath.Dir(netManagerTester.configFilePath),
	})
	var err error
	netManagerTester.manager, err = infrastructure.NewHostNetworkManager(netManagerTester.storage)
	if err != nil {
		t.Errorf("error while creating host network manager: %s", err)
	}
}

func Test_HostNetworkManager_GetList(t *testing.T) {
	links, err := netManagerTester.manager.GetList()
	if err != nil {
		t.Errorf("error getting list: %s", err.Error())
	}
	loExist := false
	for _, link := range links {
		if link.GetName() == "lo" {
			loExist = true
		}
	}
	if !loExist {
		t.Errorf("localhost not found")
	}
	err = netManagerTester.manager.SaveConfiguration()
	if err != nil {
		t.Errorf("error saving configuration: %s", err.Error())
	}
}

func Test_HostNetworkManager_CreateVlan(t *testing.T) {
	var master string
	links, err := netManagerTester.manager.GetList()
	if err != nil {
		t.Errorf("error getting list: %s", err.Error())
	}
	for _, link := range links {
		if link.GetName() != "lo" && link.GetType() != "vlan" {
			master = link.GetName()
			break
		}
	}
	netManagerTester.vlanID = 146
	netManagerTester.vlanName, err = netManagerTester.manager.CreateVlan(master, netManagerTester.vlanID)
	if err != nil {
		t.Errorf("error creating vlan: %s", err.Error())
	}
	links, err = netManagerTester.manager.GetList()
	vlanFound := false
	for _, link := range links {
		if link.GetName() == netManagerTester.vlanName {
			vlanFound = true
		}
	}
	if !vlanFound {
		t.Error("created vlan not found")
	}
	if !netManagerTester.manager.HasUnsavedChanges() {
		t.Error("after creating vlan, we must have unsaved changes")
	}
}

func Test_HostNetworkManager_CreateBridge(t *testing.T) {
	var err error
	netManagerTester.bridgeName, err = netManagerTester.manager.CreateBridge("mtest")
	if err != nil {
		t.Errorf("error creating bridge: %s", err.Error())
	}
	links, err := netManagerTester.manager.GetList()
	bridgeFound := false
	for _, link := range links {
		if link.GetName() == netManagerTester.bridgeName {
			bridgeFound = true
		}
	}
	if !bridgeFound {
		t.Error("created bridge not found")
	}
	if !netManagerTester.manager.HasUnsavedChanges() {
		t.Error("after creating bridge, we must have unsaved changes")
	}
}

func Test_HostNetworkManager_VlanSetMaster(t *testing.T) {
	err := netManagerTester.manager.SetLinkMaster(netManagerTester.vlanName, netManagerTester.bridgeName)
	if err != nil {
		t.Errorf("failed to set vlan master: %s", err.Error())
	}
	bridge, err := netManagerTester.manager.GetByName(netManagerTester.bridgeName)
	if err != nil {
		t.Errorf("error getting by name: %s", err.Error())
	}
	slaves := bridge.(domain.HostNetworkBridge).GetSlaves()
	slaveFound := false
	for _, slave := range slaves {
		if slave == netManagerTester.vlanName {
			slaveFound = true
		}
	}
	if !slaveFound {
		t.Error("the slave that was added was not found")
	}
}

func Test_HostNetworkManager_VlanAddrAdd(t *testing.T) {
	ip, ipNet, err := net.ParseCIDR("192.111.111.111/24")
	if err != nil {
		t.Errorf("parse CIDR failed: %s", err.Error())
	}
	ipNet.IP = ip
	err = netManagerTester.manager.AddrAdd(netManagerTester.vlanName, *ipNet)
	if err != nil {
		t.Errorf("failed to set the ip address: %s", err.Error())
	}
	vlan, err := netManagerTester.manager.GetByName(netManagerTester.vlanName)
	if err != nil {
		t.Errorf("error getting by name: %s", err.Error())
	}
	addresses := vlan.GetAddresses()
	addrFound := false
	for _, addr := range addresses {
		if addr.IP.Equal(ip) {
			addrFound = true
		}
	}
	if !addrFound {
		t.Error("the address that was added was not found")
	}
}

func Test_HostNetworkManager_SaveConfiguration(t *testing.T) {
	err := netManagerTester.manager.SaveConfiguration()
	if err != nil {
		t.Errorf("failed saving configuration: %s", err.Error())
	}
	conf, err := netManagerTester.storage.GetConfig()
	if err != nil {
		t.Errorf("failed to get configuration from storage: %s", err)
	}
	vlanFound := false
	for _, vlan := range conf.Vlans {
		if vlan.Name == netManagerTester.vlanName {
			vlanFound = true
		}
	}
	if !vlanFound {
		t.Error("created vlan not found in configuration file")
	}
	bridgeFound := false
	for _, bridge := range conf.Bridges {
		if bridge.Name == netManagerTester.bridgeName {
			bridgeFound = true
		}
	}
	if !bridgeFound {
		t.Error("created bridge not found in configuration file")
	}
	if netManagerTester.manager.HasUnsavedChanges() {
		t.Error("we shouldn't have unsaved changes after saving configuration")
	}
}

func Test_HostNetworkManager_VlanSetNoMaster(t *testing.T) {
	err := netManagerTester.manager.UnsetLinkMaster(netManagerTester.vlanName)
	if err != nil {
		t.Error("failed set no master for vlan")
	}
	bridge, err := netManagerTester.manager.GetByName(netManagerTester.bridgeName)
	if err != nil {
		t.Error("failed to get bridge")
	}
	slaves := bridge.(domain.HostNetworkBridge).GetSlaves()
	if len(slaves) > 0 {
		t.Error("bridge have slaves after vlan was set no master")
	}
}

func Test_HostNetworkManager_VlanAddrDelete(t *testing.T) {
	ip, ipNet, err := net.ParseCIDR("192.111.111.111/24")
	if err != nil {
		t.Errorf("parse CIDR failed: %s", err.Error())
	}
	ipNet.IP = ip
	err = netManagerTester.manager.AddrDelete(netManagerTester.vlanName, *ipNet)
	if err != nil {
		t.Errorf("failed to set the ip address: %s", err.Error())
	}
	vlan, err := netManagerTester.manager.GetByName(netManagerTester.vlanName)
	if err != nil {
		t.Errorf("error getting by name: %s", err.Error())
	}
	addresses := vlan.GetAddresses()
	for _, addr := range addresses {
		if addr.IP.Equal(ip) {
			t.Error("this address will be removed")
		}
	}
}

func Test_HostNetworkManager_ResetChangesWithDeletedVlanAddrButNotSaved(t *testing.T) {
	ip, _, err := net.ParseCIDR("192.111.111.111/24")
	if err != nil {
		t.Errorf("parse CIDR failed: %s", err.Error())
	}

	err = netManagerTester.manager.ResetChanges()
	if err != nil {
		t.Errorf("failed reset configuration to state from configuration storage: %s", err.Error())
	}
	vlan, err := netManagerTester.manager.GetByName(netManagerTester.vlanName)
	if err != nil {
		t.Errorf("error getting by name: %s", err.Error())
	}
	addresses := vlan.GetAddresses()
	addrFound := false
	for _, addr := range addresses {
		if addr.IP.Equal(ip) {
			addrFound = true
		}
	}
	if !addrFound {
		t.Error("the address that was restored was not found")
	}
}

func Test_HostNetworkManager_CheckBackupConfig(t *testing.T) {
	conf, err := netManagerTester.storage.GetBackupConfig()
	if err != nil {
		t.Errorf("get backup config failed: %s", err)
	}
	vlanFound := false
	for _, vlan := range conf.Vlans {
		if vlan.Name == netManagerTester.vlanName {
			vlanFound = true
		}
	}
	if vlanFound {
		t.Error("vlan must be deleted in backup configuration")
	}
}

func Test_HostNetworkManager_Delete(t *testing.T) {
	err := netManagerTester.manager.DeleteLinkByName(netManagerTester.vlanName)
	if err != nil {
		t.Errorf("failed deleting vlan: %s", err.Error())
	}
	links, err := netManagerTester.manager.GetList()
	vlanFound := false
	for _, link := range links {
		if link.GetName() == netManagerTester.vlanName {
			vlanFound = true
		}
	}
	if vlanFound {
		t.Error("deleted vlan was found")
	}
}

func Test_HostNetworkManager_ResetChanges(t *testing.T) {
	err := netManagerTester.manager.ResetChanges()
	if err != nil {
		t.Errorf("failed reset configuration to state from configuration storage: %s", err.Error())
	}
	links, err := netManagerTester.manager.GetList()
	vlanFound := false
	for _, link := range links {
		if link.GetName() == netManagerTester.vlanName {
			vlanFound = true
		}
	}
	if !vlanFound {
		t.Error("vlan not found at host")
	}
}

func Test_HostNetworkManager_RestoreFromBackup(t *testing.T) {
	err := netManagerTester.manager.RestoreFromBackup()
	if err != nil {
		t.Errorf("failed restore configuration from backup: %s", err.Error())
	}
	links, err := netManagerTester.manager.GetList()
	vlanExist := false
	for _, link := range links {
		if link.GetType() == "vlan" && link.GetName() == netManagerTester.vlanName {
			vlanExist = true
		}
	}
	if vlanExist {
		t.Error("vlan was found on host")
	}
	if !netManagerTester.manager.HasUnsavedChanges() {
		t.Error("after restore from backup, we must have unsaved changes")
	}
}

func Test_HostNetworkManager_CleaningAfterTests(t *testing.T) {
	err := os.Remove(netManagerTester.configFilePath)
	if err != nil {
		t.Errorf("remove network config file failed:  %q", err)
	}
	err = os.Remove(netManagerTester.configFilePath + ".back")
	if err != nil {
		t.Errorf("remove network backup config file failed:  %q", err)
	}
}
