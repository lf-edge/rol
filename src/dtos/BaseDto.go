package dtos

import (
	"time"
)

//BaseDto base dto type
type BaseDto[IDType comparable] struct {
	//	ID - unique identifier
	ID IDType
	//	CreatedAt - entity create time
	CreatedAt time.Time
	//	UpdatedAt - entity update time
	UpdatedAt *time.Time
}
