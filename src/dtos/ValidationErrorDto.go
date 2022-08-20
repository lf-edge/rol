package dtos

//ValidationErrorElemDto element of ValidationErrorDto
type ValidationErrorElemDto struct {
	//Source for the field, for example query or body
	Source string
	//Field with validation error
	Field string
	//Value for Key
	Error string
}

//ValidationErrorDto validation error dto
type ValidationErrorDto struct {
	Message string
	Errors  []ValidationErrorElemDto
}
