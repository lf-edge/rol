package errors

import (
	"fmt"
	"runtime"
)

// New creates an error
//Params
//	ErrorType - error type
//	msg - error message
//Return
//	error - new generated error
func New(errorType ErrorType, msg string) error {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	return customError{
		stack:       string(stackSlice[0:s]),
		errorType:   errorType,
		msg:         msg,
		err:         nil,
		contextInfo: make(map[string]string),
	}
}

//Newf creates an error
//Params
//	ErrorType - error type
//	msg - format for error message
//	args - arguments for message formatting
//Return
//	error - new generated error
func Newf(errorType ErrorType, msg string, args ...interface{}) error {
	return New(errorType, fmt.Sprintf(msg, args...))
}

//Wrap wraps an error
//Params
//	error - error to wrapping
//	ErrorType - error type
//	msg - error message
//Return
//	error - new generated wrapped error
func Wrap(err error, errorType ErrorType, msg string) error {
	stackSlice := make([]byte, 512)
	s := runtime.Stack(stackSlice, false)
	stackTrace := string(stackSlice[0:s])

	customErr := customError{
		stack:       "",
		errorType:   errorType,
		msg:         fmt.Sprintf("%s: %s", msg, err.Error()),
		err:         err,
		contextInfo: make(map[string]string),
	}

	// Copy already existed stack trace if wrapped error is custom error
	if customErrForWrap, ok := err.(customError); ok {
		customErr.stack = customErrForWrap.stack
	} else {
		customErr.stack = stackTrace
	}

	return customErr
}

//Wrapf wraps an error
//Params
//	error - error to wrapping
//	ErrorType - error type
//	msg - format for error message
//	args - arguments for message formatting
//Return
//	error - new generated wrapped error
func Wrapf(err error, errorType ErrorType, msg string, args ...interface{}) error {
	return Wrap(err, errorType, fmt.Sprintf(msg, args...))
}

//Unwrap gives the wrapped error
//Params
//	error - error to unwrap
//Return
//	error - unwrapped error or nil
func Unwrap(err error) error {
	if customErr, ok := err.(customError); ok {
		return customErr.err
	}
	return nil
}
