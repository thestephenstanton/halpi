package errors

// HapiError is a custom error that helps deliver status codes from deeper in
// code to your application layer
type HapiError struct {
	Err       error
	ErrorType ErrorType

	// Message is the original message that you pass to create the new hapi error
	// example errors.BadRequest.Wrap(err, "this would be the message")
	Message string
}

// Error returns the error string of a HapiError.
func (e HapiError) Error() string {
	return e.Err.Error()
}

// GetStatusCode gets the status code for the HapiError.
func (e HapiError) GetStatusCode() int {
	return getStatusCode(e.ErrorType)
}

// GetMessage gets the Message of the HapiError.
func (e HapiError) GetMessage() string {
	return e.Message
}

// SetMessage sets the Message and returns new error with new Message set.
func (e HapiError) SetMessage(message string) HapiError {
	e.Message = message

	return e
}

// SetMessage will set the Message of a HapiError so that you
// can return a detailed message for the client when responding.
// If err is not of type HapiError, it will be converted to a NoType
// HapiError and have the message set.
func SetMessage(err error, message string) HapiError {
	hapiError := CastToHapiError(err)

	hapiError.Message = message

	return hapiError
}

// CastToHapiError turns normal error into HapiError. If already a HapiError, this
// will have no effect. If it is not, then this will return a NoType HapiError.
func CastToHapiError(err error) HapiError {
	var hapiError HapiError
	ok := As(err, &hapiError)
	if !ok {
		return HapiError{
			ErrorType: NoType,
			Err:       err,
		}
	}

	return hapiError
}
