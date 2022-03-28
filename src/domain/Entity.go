package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Entity base entity structure
type Entity struct {
	//	ID - entity identifier
	ID uuid.UUID `gorm:"primary_key;"`
	//	CreatedAt - time when the entity was created
	CreatedAt time.Time `gorm:"index"`
	//	UpdatedAt - time when the entity was updated
	UpdatedAt time.Time `gorm:"index"`
	//	DeletedAt - time when the entity was deleted
	DeletedAt time.Time `gorm:"index"`
}

//GetID gets the id of the entity
func (ent Entity) GetID() uuid.UUID {
	return ent.ID
}

//SetDeleted set the entity DeletedAt field at time.Now()
func (ent *Entity) SetDeleted() {
	ent.DeletedAt = time.Now()
}

// BeforeCreate will set a UUID rather than numeric ID.
func (ent *Entity) BeforeCreate(_ *gorm.DB) (err error) {
	ent.ID = uuid.New()
	return
}
