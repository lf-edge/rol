// Package domain stores the main structures of the program
package domain

//TFTPConfig TFTP config entity
type TFTPConfig struct {
	EntityUUID
	//Address TFTP server IP address
	Address string
	//Port TFTP server port
	Port string
	//Enabled TFTP server startup status
	Enabled bool
}
