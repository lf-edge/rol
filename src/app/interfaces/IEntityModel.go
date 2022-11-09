package interfaces

import (
	"time"
)

// IEntityModel default interface that must be implemented by each entity
type IEntityModel[IDType comparable] interface {
	//GetID of the entity
	//
	//Return
	//	IDType - entity id
	GetID() IDType
	//GetCreatedAt time
	//
	//Return
	//	uuid.UUID - create time
	GetCreatedAt() time.Time
	//GetUpdatedAt time
	//
	//Return
	//	*time.Time - update time
	GetUpdatedAt() *time.Time
	//GetDeletedAt time
	//
	//Return
	//	time.Time - update time
	GetDeletedAt() *time.Time
}
