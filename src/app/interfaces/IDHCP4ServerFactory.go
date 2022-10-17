package interfaces

import "rol/domain"

//IDHCP4ServerFactory interface for DHCP v4 server fabric implementations
type IDHCP4ServerFactory interface {
	//Create DHCP v4 server with config
	//
	//Params:
	//	config - dhcp v4 config
	//Return:
	//  IDHCP4Server - dhcp v4 server
	//	error - if an error occurred, otherwise nil
	Create(config domain.DHCP4Config) (IDHCP4Server, error)
}
