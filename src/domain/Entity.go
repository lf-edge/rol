package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Entity base entity structure
type Entity struct {
	//	ID - entity identifier
	ID uuid.UUID `gorm:"primary_key; type:varchar(36)"`
	//	CreatedAt - time when the entity was created
	CreatedAt time.Time `gorm:"index"`
	//	UpdatedAt - time when the entity was updated
	UpdatedAt time.Time `gorm:"index"`
	//	DeletedAt - time when the entity was deleted
	DeletedAt gorm.DeletedAt
}

//GetID gets the id of the entity
func (e Entity) GetID() uuid.UUID {
	return e.ID
}

// BeforeCreate will set a UUID rather than numeric ID.
func (e *Entity) BeforeCreate(_ *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}
