package dtos

import "time"

type BaseDto struct {
	Id        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
