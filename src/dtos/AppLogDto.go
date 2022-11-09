package dtos

import (
	"github.com/google/uuid"
)

//AppLogDto dto for app log
type AppLogDto struct {
	BaseDto[uuid.UUID]
	ActionID uuid.UUID `gorm:"index"`
	//	Level - level of the log
	Level string
	//	Source - method from which the log was obtained
	Source string
	//	Message - log message
	Message string
}
