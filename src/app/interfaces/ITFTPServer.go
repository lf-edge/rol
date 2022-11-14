package interfaces

import "rol/domain"

//ITFTPServer define interface for TFTP server implementation
type ITFTPServer interface {
	//ReloadConfig for TFTP server
	ReloadConfig(config domain.TFTPConfig) error
	//ReloadPaths for TFTP server on fly (without stop/start)
	ReloadPaths(paths []domain.TFTPPathRatio) error
	//Start TFTP server
	Start() error
	//Stop TFTP server
	Stop()
	//GetState of TFTP server
	GetState() domain.TFTPServerState
}
