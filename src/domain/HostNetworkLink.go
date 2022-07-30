package domain

import "net"

//HostNetworkLink is a struct for base network link
type HostNetworkLink struct {
	//Name link name
	Name string
	//Type link type
	Type string
	//Addresses slice of link addresses net.IPNet
	Addresses []net.IPNet `yaml:"addresses"`
}

//GetName gets host interface name
func (h HostNetworkLink) GetName() string {
	return h.Name
}

//GetType gets host interface type
func (h HostNetworkLink) GetType() string {
	return h.Type
}

//GetAddresses gets host interface addresses
func (h HostNetworkLink) GetAddresses() []net.IPNet {
	return h.Addresses
}
