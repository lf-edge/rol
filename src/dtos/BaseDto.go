package dtos

import (
	"time"

	"github.com/google/uuid"
)

//BaseDto base dto type
type BaseDto struct {
	//	Id - unique identifier
	Id uuid.UUID
	//	CreatedAt - entity create time
	CreatedAt time.Time
	//	UpdatedAt - entity update time
	UpdatedAt time.Time
}
