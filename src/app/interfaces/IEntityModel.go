package interfaces

import "github.com/google/uuid"

// IEntityModel default interface that must be implemented by each entity
type IEntityModel interface {
	//GetID
	//	Gets entity id
	//Return
	//	uuid.UUID - entity id
	GetID() uuid.UUID
}
