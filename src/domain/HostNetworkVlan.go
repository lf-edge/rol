package domain

//HostNetworkVlan is a struct for network vlan
type HostNetworkVlan struct {
	HostNetworkLink
	//VlanID vlan ID
	VlanID int
	//Parent name of parent network interface
	Parent string
}
