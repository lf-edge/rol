package dtos

import (
	"github.com/google/uuid"
)

//AppLogDto dto for app log
type AppLogDto struct {
	BaseDto
	ActionID uuid.UUID `gorm:"index"`
	//	Level - level of the log
	Level string
	//	Source - method from which the log was obtained
	Source string
	//	Message - log message
	Message string
}

//Validate validates dto fields
//Return
//	error - if error occurs return error, otherwise nil
func (hld AppLogDto) Validate() error {
	return nil
}
