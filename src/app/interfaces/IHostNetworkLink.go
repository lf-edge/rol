package interfaces

import "net"

//IHostNetworkLink is an interface for base network link
type IHostNetworkLink interface {
	//GetName gets host interface name
	GetName() string
	//GetType gets host interface type
	GetType() string
	//GetAddresses gets host interface addresses
	GetAddresses() []net.IPNet
}
