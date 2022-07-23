package errors

//AddErrorContext adds a context to an error
//Params
//	error - error for adding context
//	key - key for context info
//	value - value for context info
//Return
//	error - error with added context
func AddErrorContext(err error, key, value string) error {
	if customErr, ok := err.(customError); ok {
		customErr.context[key] = value
		return customErr
	}
	context := make(map[string]string)
	context[key] = value
	return customError{errorType: NoType, msg: err.Error(), wrapped: nil, context: context}
}

//GetErrorContext returns the error context
//Params
//	error - error for getting context
//Return
//	context - key and value context map
func GetErrorContext(err error) map[string]string {
	if customErr, ok := err.(customError); ok {
		return customErr.context
	}
	return make(map[string]string)
}

//GetType returns the error type
//Params
//	error - error for getting type
//Return
//	ErrorType - error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}
	return NoType
}

//As determines whether there was an error with the specified type in the error stack
//Params
//	error - error
//	ErrorType - error type
//Return
//	bool - true if error type exist in the wrapped errors
func As(err error, errorType ErrorType) bool {
	if customErr, ok := err.(customError); ok {
		if customErr.errorType == errorType {
			return true
		}

		unwrappedErr := Unwrap(err)
		if unwrappedErr != nil {
			return As(unwrappedErr, errorType)
		}
	}
	return false
}

//GetCallerFile Get Caller Line from stack trace
//Params
//	error - error
//Return
//	string - full path to the caller file or empty string if stack is not exist
func GetCallerFile(err error) string {
	if customErr, ok := err.(customError); ok {
		return customErr.file
	}
	return ""
}

//GetCallerLine Get Caller Line from stack trace
//Params
//	error - error
//Return
//	int - -1, if caller line is not exist in error
func GetCallerLine(err error) int {
	if customErr, ok := err.(customError); ok {
		return customErr.line
	}
	return -1
}
