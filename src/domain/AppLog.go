package domain

import (
	"github.com/google/uuid"
)

//AppLog App log entity
type AppLog struct {
	//EntityUUID - nested base entity where ID type is uuid.UUID
	EntityUUID
	//	ActionID - http request id, empty if the log is initiated inside the application
	ActionID uuid.UUID `gorm:"index"`
	//	Level - level of the log
	Level string
	//	Source - method from which the log was obtained
	Source string
	//	Message - log message
	Message string
}

//GetID gets the id of the log
func (a AppLog) GetID() uuid.UUID {
	return a.ID
}
