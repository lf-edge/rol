package errors

import (
	"fmt"
	"runtime"
)

func newErr(skipCaller int, errorType ErrorType, msg string) error {
	_, file, line, _ := runtime.Caller(skipCaller)
	return customError{
		file:      file,
		line:      line,
		errorType: errorType,
		msg:       msg,
		wrapped:   nil,
		context:   make(map[string]string),
	}
}

func wrapErr(err error, skipCaller int, errorType ErrorType, msg string) error {
	_, file, line, _ := runtime.Caller(skipCaller)
	customErr := customError{
		file:      "",
		line:      0,
		errorType: errorType,
		msg:       fmt.Sprintf("%s: %s", msg, err.Error()),
		wrapped:   err,
		context:   make(map[string]string),
	}

	// Copy already existed stack trace if wrapped error is custom error
	if customErrForWrap, ok := err.(customError); ok {
		customErr.file = customErrForWrap.file
		customErr.line = customErrForWrap.line
	} else {
		customErr.file = file
		customErr.line = line
	}

	return customErr
}

// New creates an error
//Params
//	msg - error message
//Return
//	error - new generated error
func New(msg string) error {
	return newErr(2, NoType, msg)
}

//Newf creates an error
//Params
//	msg - format for error message
//	args - arguments for message formatting
//Return
//	error - new generated error
func Newf(msg string, args ...interface{}) error {
	return newErr(2, NoType, fmt.Sprintf(msg, args...))
}

//Wrap wraps an error
//Params
//	error - error to wrapping
//	msg - error message
//Return
//	error - new generated wrapped error
func Wrap(err error, msg string) error {
	return wrapErr(err, 2, NoType, msg)
}

//Wrapf wraps an error
//Params
//	error - error to wrapping
//	msg - format for error message
//	args - arguments for message formatting
//Return
//	error - new generated wrapped error
func Wrapf(err error, msg string, args ...interface{}) error {
	return wrapErr(err, 2, NoType, fmt.Sprintf(msg, args...))
}

//Unwrap gives the wrapped error
//Params
//	error - error to unwrap
//Return
//	error - unwrapped error or nil
func Unwrap(err error) error {
	if customErr, ok := err.(customError); ok {
		return customErr.wrapped
	}
	return nil
}
