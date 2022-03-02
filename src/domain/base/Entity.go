package base

import "time"

type Entity struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
	DeletedAt time.Time `gorm:"index"`
}

func (ent *Entity) GetId() uint {
	return ent.ID
}

func (ent *Entity) SetDeleteTimeNow() {
	ent.DeletedAt = time.Now()
}
