package domain

//HostNetworkConfig is a struct for yaml configuration file
type HostNetworkConfig struct {
	//Devices slice of HostNetworkDevice
	Devices []HostNetworkDevice
	//Vlans slice of HostNetworkVlan
	Vlans []HostNetworkVlan
}
