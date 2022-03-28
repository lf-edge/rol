package interfaces

// IEntityDtoModel default interface that must be implemented by each entity
type IEntityDtoModel interface {
	//Validate
	//	Validates dto fields
	//Return
	//	error - if error occurs return error, otherwise nil
	Validate() error
}
