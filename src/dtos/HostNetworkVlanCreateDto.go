package dtos

//HostNetworkVlanCreateDto create dto for HostNetworkVlan entity
type HostNetworkVlanCreateDto struct {
	//VlanID vlan id
	VlanID int
	//Parent interface name
	Parent string
	//Addresses list
	Addresses []string
}
