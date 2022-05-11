package dtos

import (
	"time"

	"github.com/google/uuid"
)

//BaseDto base dto type
type BaseDto struct {
	//	ID - unique identifier
	ID uuid.UUID
	//	CreatedAt - entity create time
	CreatedAt time.Time
	//	UpdatedAt - entity update time
	UpdatedAt time.Time
}
