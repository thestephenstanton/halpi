package errors

// HapiError is a custom error that helps deliver status codes from deeper in
// code to your application layer
type HapiError struct {
	Err       error
	ErrorType ErrorType
	Message   string
}

// Error returns the error string of a HapiError.
func (hapiErr HapiError) Error() string {
	return hapiErr.Err.Error()
}

// GetStatusCode gets the status code for the HapiError.
func (hapiErr HapiError) GetStatusCode() int {
	return getStatusCode(hapiErr.ErrorType)
}

// GetMessage gets the Message of the HapiError.
func (hapiErr HapiError) GetMessage() string {
	return hapiErr.Message
}

// SetMessage will set the Message of a HapiError so that you
// can return a detailed message for the client when responding.
// If err is not of type HapiError, it will be converted to a NoType
// HapiError and have the message set.
func SetMessage(err error, message string) error {
	hapiError := castToHapiError(err)

	hapiError.Message = message

	return hapiError
}

// Cast turns normal error into HapiError. If already a HapiError, this
// will have no effect. If it is not, then this will return a NoType HapiError.
func Cast(err error) error {
	return castToHapiError(err)
}

func castToHapiError(err error) HapiError {
	hapiErr, ok := err.(HapiError)
	if !ok {
		return HapiError{
			ErrorType: NoType,
			Err:       err,
		}
	}

	return hapiErr
}
