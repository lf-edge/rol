package interfaces

import "rol/domain"

//IDHCP4Server interface for DHCP v4 server implementations
type IDHCP4Server interface {
	//ReloadConfiguration DHCP v4 server from config
	ReloadConfiguration(dhcp4config domain.DHCP4Config) error
	//Start DHCP v4 server
	Start() error
	//Stop DHCP v4 server
	Stop()
	//GetState of DHCP v4 server
	GetState() domain.DHCPServerState
}
