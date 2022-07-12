package errors

//ErrorType is a custom error type
type ErrorType uint

const (
	//NoType error without type
	NoType = ErrorType(iota)
	//Internal error type
	Internal
	//Validation error type
	Validation
	//NotFound error type
	NotFound
	//AlreadyExist error type
	AlreadyExist
)

type customError struct {
	stack       string
	msg         string
	err         error
	errorType   ErrorType
	contextInfo map[string]string
}

// Error returns the error message
//Return
//	string - error message
func (c customError) Error() string {
	return c.msg
}

//New creates a new error
//Params
//	msg - error message
//Return
//	error - new generated error
func (e ErrorType) New(msg string) error {
	return New(e, msg)
}

//Newf creates a new error with formatted message
//Params
//	msg - format for error message
//	args - arguments for message formatting
//Return
//	error - new generated error
func (e ErrorType) Newf(msg string, args ...interface{}) error {
	return Newf(e, msg, args...)
}

//Wrapf wraps an error
//Params
//	error - error to wrapping
//	msg - format for error message
//	args - arguments for message formatting
//Return
//	error - new generated wrapped error
func (e ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return Wrapf(err, e, msg, args...)
}

//Wrap wraps an error
//Params
//	error - error to wrapping
//	msg - error message
//Return
//	error - new generated wrapped error
func (e ErrorType) Wrap(err error, msg string) error {
	return Wrap(err, e, msg)
}
