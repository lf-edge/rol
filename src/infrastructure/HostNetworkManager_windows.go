package infrastructure

import (
	"net"
	"rol/app/interfaces"
)

//HostNetworkManager is a struct for network manager
type HostNetworkManager struct {
	configStorage interfaces.IHostNetworkConfigStorage
}

//NewHostNetworkManager constructor for HostNetworkManager
func NewHostNetworkManager(configStorage interfaces.IHostNetworkConfigStorage) (interfaces.IHostNetworkManager, error) {
	hostNetworkManager := &HostNetworkManager{
		configStorage: configStorage,
	}
	return hostNetworkManager, nil
}

//GetList gets list of host network interfaces
//
//Return:
//	[]interfaces.IHostNetworkLink - list of interfaces
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) GetList() ([]interfaces.IHostNetworkLink, error) {
	panic("not implemented")
}

//GetByName gets host network interface by its name
//
//Params:
//	name - interface name
//Return:
//	interfaces.IHostNetworkLink - interfaces
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) GetByName(_ string) (interfaces.IHostNetworkLink, error) {
	panic("not implemented")
}

//CreateVlan creates vlan on host
//
//Params:
//	master - name of the master network interface
//	vlanID - ID to be set for vlan
//Return:
//	string - new vlan name that will be rol.{master}.{vlanID}
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) CreateVlan(_ string, _ int) (string, error) {
	panic("not implemented")
}

//CreateBridge creates bridge on host
//
//Params:
//	name - new bridge name
//Return:
//	string - new bridge name that will be rol.br.{name}
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) CreateBridge(_ string) (string, error) {
	panic("not implemented")
}

//SetLinkUp enables the link
//
//Params:
//	linkName - name of the link
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) SetLinkUp(_ string) error {
	panic("not implemented")
}

//SetLinkMaster set master for link
//
//Params:
//	slaveName - name of link that will be slave
//	masterName - name of link that will be master for the slave
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) SetLinkMaster(_, _ string) error {
	panic("not implemented")
}

//UnsetLinkMaster removes the master of the link
//
//Params:
//	linkName - name of the link
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) UnsetLinkMaster(_ string) error {
	panic("not implemented")
}

//DeleteLinkByName deletes interface on host by its name
//
//Params:
//	name - interface name
//Return
//	error - if an error occurs, otherwise nil

//DeleteLinkByName deletes interface on host by its name
//
//Params:
//	name - interface name
//Return
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) DeleteLinkByName(_ string) error {
	panic("not implemented")
}

//AddrAdd Add new ip address for network interface
//
//Params:
//	linkName - name of the interface
//	addr - ip address with mask net.IPNet
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) AddrAdd(_ string, _ net.IPNet) error {
	panic("not implemented")
}

//AddrDelete Delete ip address from network interface
//
//Params:
//	linkName - name of the interface
//	addr - ip address with mask net.IPNet
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) AddrDelete(_ string, _ net.IPNet) error {
	panic("not implemented")
}

//SaveConfiguration save current host network configuration to the configuration storage
//Save previous config file to .back file
//
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) SaveConfiguration() error {
	panic("not implemented")
}

//RestoreFromBackup restore and apply host network configuration from backup
//
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) RestoreFromBackup() error {
	panic("not implemented")
}

//ResetChanges Reset all applied changes to state from saved configuration
//
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkManager) ResetChanges() error {
	panic("not implemented")
}

//HasUnsavedChanges Gets a flag about unsaved changes
//
//Return:
//	bool - if unsaved changes exist - true, otherwise false
func (h *HostNetworkManager) HasUnsavedChanges() bool {
	panic("not implemented")
}
