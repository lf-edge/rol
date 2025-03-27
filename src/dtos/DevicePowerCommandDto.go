package dtos

import (
	"github.com/google/uuid"
)

type DevicePowerCommandDto struct {
	BaseDto[uuid.UUID]
	Command string
}
