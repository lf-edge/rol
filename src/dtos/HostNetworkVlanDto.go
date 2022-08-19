package dtos

//HostNetworkVlanDto dto for HostNetworkVlan entity
type HostNetworkVlanDto struct {
	//Name interface full name
	Name string
	//Addresses list
	Addresses []string
	//VlanID id of the vlan
	VlanID int
	//Master interface name
	Master string
}
