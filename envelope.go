package hapi

// ResponseEnvelope wraps the data returned to the client in an envelope
type ResponseEnvelope struct {
	StatusCode int           `json:"status,omitifempty"`
	Data       interface{}   `json:"data,omitifempty"`
	Error      ErrorEnvelope `json:"error,omitifempty"`
}

// ErrorEnvelope wraps the error returned to the client in an envelope that is
// a part of the ResponseEnvelope
type ErrorEnvelope struct {
	Message string `json:"message,omitifempty"`
}

// NewResponseEnvelope creates a new ResponseEnvelope
func NewResponseEnvelope(statusCode int, data interface{}) ResponseEnvelope {
	return ResponseEnvelope{
		StatusCode: statusCode,
		Data:       data,
	}
}

type hapiError interface {
	GetStatusCode() int
	GetMessage() string
}

// NewErrorEnvelope creates a ResponseEnvelope that contains an error
func NewErrorEnvelope(statusCode int, err error) ResponseEnvelope {
	envelope := ResponseEnvelope{
		StatusCode: statusCode,
	}

	// Check if it is a compatible hapiError
	hapiErr, ok := err.(hapiError)
	if !ok {
		envelope.Error.Message = err.Error()
		return envelope
	}

	envelope.StatusCode = hapiErr.GetStatusCode()
	envelope.Error.Message = hapiErr.GetMessage()

	return envelope
}
