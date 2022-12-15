package dtos

import "github.com/google/uuid"

type ProjectDto struct {
	BaseDto[uuid.UUID]
	ProjectBaseDto
	BridgeName   string
	DHCPServerID uuid.UUID
	TFTPServerID uuid.UUID
}
