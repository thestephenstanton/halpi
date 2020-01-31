package errors

type hapiError struct {
	err       error
	errorType ErrorType
	message   string
}

// Error returns the message of a hapiError
func (hapiError hapiError) Error() string {
	return hapiError.err.Error()
}

func (hapiError hapiError) GetStatusCode() int {
	return getStatusCode(hapiError.errorType)
}

func (hapiError hapiError) GetMessage() string {
	return hapiError.message
}

// SetMessage will set the message of a hapiError
// so that the envelope can be properly set when returning a response
// if err is not of type hapiError. Returns the error with the
// message set and a bool if it was set or not. Bool will be false
// if the error passed in is not of type hapiError
func SetMessage(err error, message string) (error, bool) {
	hapiError, ok := castToHapiError(err)
	if !ok {
		return err, false
	}

	hapiError.message = message

	return hapiError, true
}

func castToHapiError(err error) (hapiError, bool) {
	customErr, ok := err.(hapiError)
	if !ok {
		return hapiError{}, false
	}

	return customErr, true
}
