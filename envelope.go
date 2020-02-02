package hapi

// ResponseEnvelope wraps the data returned to the client in an envelope
type ResponseEnvelope struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

// ErrorEnvelope wraps the error returned to the client in an envelope that is
// a part of the ResponseEnvelope
type ErrorEnvelope struct {
	Message string `json:"message,omitempty"`
}

// NewResponseEnvelope creates a new ResponseEnvelope
func NewResponseEnvelope(statusCode int, data interface{}) ResponseEnvelope {
	return ResponseEnvelope{
		StatusCode: statusCode,
		Data:       data,
	}
}

type hapiError interface {
	error
	GetStatusCode() int
	GetMessage() string
}

// NewErrorEnvelope creates a ResponseEnvelope that contains an error.
func NewErrorEnvelope(statusCode int, payload interface{}) ResponseEnvelope {
	envelope := ResponseEnvelope{
		StatusCode: statusCode,
	}

	// Check if it is a hapiError.
	hapiErr, ok := payload.(hapiError)
	if ok {
		envelope.StatusCode = hapiErr.GetStatusCode()

		if Config.UseHapiEnvelopes {
			envelope.Error = ErrorEnvelope{
				Message: hapiErr.GetMessage(),
			}
		} else {
			envelope.Error = hapiErr.GetMessage()
		}

		return envelope
	}

	// Check if it is a regular error.
	regularErr, ok := payload.(error)
	if ok {
		if Config.UseHapiEnvelopes {
			envelope.Error = ErrorEnvelope{
				Message: regularErr.Error(),
			}
		} else {
			envelope.Error = regularErr.Error()
		}

		return envelope
	}

	// Otherwise, just stick it in there.
	envelope.Error = payload

	return envelope
}
