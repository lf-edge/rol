package domain

import "github.com/google/uuid"

//Project entity of a project
type Project struct {
	EntityUUID
	//Name project name
	Name string
	//BridgeName name of a bridge that related to project
	BridgeName string
	//Subnet project subnet
	Subnet string
	//DHCPServerID project DHCP server ID
	DHCPServerID uuid.UUID `gorm:"type:varchar(36);index"`
	//TFTPServerID project TFTP server ID
	TFTPServerID uuid.UUID `gorm:"type:varchar(36);index"`
}
