package interfaces

import "net"

//IHostNetworkManager is an interface for network manager
type IHostNetworkManager interface {
	//GetList gets list of host network interfaces
	//
	//Return:
	//	[]interfaces.IHostNetworkLink - list of interfaces
	//	error - if an error occurs, otherwise nil
	GetList() ([]IHostNetworkLink, error)
	//GetByName gets host network interface by its name
	//
	//Params:
	//	name - interface name
	//Return:
	//	interfaces.IHostNetworkLink - host network interface
	//	error - if an error occurs, otherwise nil
	GetByName(name string) (IHostNetworkLink, error)
	//CreateVlan creates vlan on host
	//
	//Params:
	//	master - name of the master network interface
	//	vlanID - ID to be set for vlan
	//Return:
	//	string - new vlan name that will be {master}.{vlanID}
	//	error - if an error occurs, otherwise nil
	CreateVlan(master string, vlanID int) (string, error)
	//DeleteLinkByName deletes interface on host by its name
	//
	//Params:
	//	name - interface name
	//Return
	//	error - if an error occurs, otherwise nil
	DeleteLinkByName(name string) error
	//AddrAdd Add new ip address for network interface
	//
	//Params:
	//	linkName - name of the interface
	//	addr - ip address with mask net.IPNet
	//Return:
	//	error - if an error occurs, otherwise nil
	AddrAdd(linkName string, addr net.IPNet) error
	//AddrDelete Delete ip address from network interface
	//
	//Params:
	//	linkName - name of the interface
	//	addr - ip address with mask net.IPNet
	//Return:
	//	error - if an error occurs, otherwise nil
	AddrDelete(linkName string, addr net.IPNet) error
	//SaveConfiguration save current host network configuration to the configuration storage
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	SaveConfiguration() error
	//RestoreFromBackup restore and apply host network configuration from backup configuration
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	RestoreFromBackup() error
	//ResetChanges Load and apply host network configuration from saved configuration
	//
	//Return:
	//	error - if an error occurs, otherwise nil
	ResetChanges() error
	//HasUnsavedChanges Gets a flag about unsaved changes
	//
	//Return:
	//	bool - if unsaved changes exist - true, otherwise false
	HasUnsavedChanges() bool
}
