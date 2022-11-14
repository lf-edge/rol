package interfaces

import "rol/domain"

//ITFTPServerFactory interface for TFTP server fabric implementations
type ITFTPServerFactory interface {
	//Create TFTP server
	//
	//Params:
	//	config - tftp server config
	//Return:
	//  ITFTPServer - tftp server
	//	error - if an error occurred, otherwise nil
	Create(config domain.TFTPConfig) (ITFTPServer, error)
}
