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
		customErr.contextInfo[key] = value
		return customErr
	}
	context := make(map[string]string)
	context[key] = value
	return customError{errorType: NoType, msg: err.Error(), err: nil, contextInfo: context}
}

//GetErrorContext returns the error context
//Params
//	error - error for getting context
//Return
//	context - key and value context map
func GetErrorContext(err error) map[string]string {
	if customErr, ok := err.(customError); ok {
		return customErr.contextInfo
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
