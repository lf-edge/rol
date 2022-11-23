// Package domain stores the main structures of the program
package domain

//HostNetworkConfig is a struct for yaml configuration file
type HostNetworkConfig struct {
	//Devices slice of HostNetworkDevice
	Devices []HostNetworkDevice
	//Vlans slice of HostNetworkVlan
	Vlans []HostNetworkVlan
	//Bridges slice of HostNetworkBridge
	Bridges []HostNetworkBridge
	//TrafficRules netfilter traffic rules struct
	TrafficRules TrafficRules
}

//TrafficRules is a struct for netfilter rules separated by tables
type TrafficRules struct {
	//Filter 'filter' table rules
	Filter []HostNetworkTrafficRule
	//NAT 'nat' table rules
	NAT []HostNetworkTrafficRule
	//Mangle 'mangle' table rules
	Mangle []HostNetworkTrafficRule
	//Raw 'raw' table rules
	Raw []HostNetworkTrafficRule
	//Security 'security' table rules
	Security []HostNetworkTrafficRule
}
