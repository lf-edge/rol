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

// IEntityModelDeletedAt default interface that must be implemented by each entity
type IEntityModelDeletedAt interface {
	//SetDeleted set the entity DeletedAt field at time.Now()
	SetDeleted()
}
