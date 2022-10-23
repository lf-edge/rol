// Package dtos stores all data transfer objects
package dtos

//TFTPServerBaseDto TFTP server base dto
type TFTPServerBaseDto struct {
	//Address TFTP server IP address
	Address string
	//Port TFTP server port
	Port string
	//Enabled TFTP server startup status
	Enabled bool
}
