package domain

//HostNetworkVlan is a struct for network vlan
type HostNetworkVlan struct {
	HostNetworkLink
	//VlanID vlan ID
	VlanID int
	//Master name of master network interface
	Master string
}
