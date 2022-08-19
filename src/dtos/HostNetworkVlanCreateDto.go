package dtos

//HostNetworkVlanCreateDto create dto for HostNetworkVlan entity
type HostNetworkVlanCreateDto struct {
	//VlanID vlan id
	VlanID int
	//Master interface name
	Master string
	//Addresses list
	Addresses []string
}
