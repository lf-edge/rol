// Package dtos stores all data transfer objects
package dtos

import "github.com/google/uuid"

//ProjectDto project dto
type ProjectDto struct {
	BaseDto[uuid.UUID]
	ProjectBaseDto
	//BridgeName name of project bridge
	BridgeName string
	//DHCPServerID ip of project DHCP server
	DHCPServerID uuid.UUID
	//TFTPServerID ip of project TFTP server
	TFTPServerID uuid.UUID
}
