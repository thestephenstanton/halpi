package hapi

// ErrorResponse is a standard error response that has an Error and a RawError
type ErrorResponse struct {
	ErrorMessage string `json:"error"`
	RawError     string `json:"rawError,omitempty"`
}

// NewErrorResponse creates new ErrorResponse with an error message.NewErrorResponse.
// To set the raw error, use NewErrorResponse("error message").SetRawError("raw error")
func NewErrorResponse(errorMessage string) ErrorResponse {
	return ErrorResponse{
		ErrorMessage: errorMessage,
	}
}

// SetRawError sets the RawError that is useful for local development
func (e ErrorResponse) SetRawError(rawError string) ErrorResponse {
	e.RawError = rawError

	return e
}

// Error adhears to error interface to get the raw error
func (e ErrorResponse) Error() string {
	return e.RawError
}
