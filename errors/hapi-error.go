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
// if err is not of type hapiError. If the error passed in is not a hapi
// error, it will be converted to a NoType hapi error and have the message set
func SetMessage(err error, message string) error {
	hapiError := castToHapiError(err)

	hapiError.message = message

	return hapiError
}

func castToHapiError(err error) hapiError {
	customErr, ok := err.(hapiError)
	if !ok {
		return hapiError{
			errorType: NoType,
			err:       err,
		}
	}

	return customErr
}
