package hapi

import (
	"encoding/json"
	"net/http"

	"github.com/thestephenstanton/hapi/errors"
)

type hapiError interface {
	error
	GetStatusCode() int
	GetMessage() string
}

type errorResponse struct {
	Error    string `json:"error"`
	RawError string `json:"rawError,omitempty"`
}

// Respond will marshal and return the payload to the client with a given status code.
func Respond(w http.ResponseWriter, statusCode int, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return errors.InternalServerError.Wrap(err, "failed to marshal payload")
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	w.Write(bytes)

	return nil
}

// RespondError will find if the error is a hapiError and if it is, get the message and set it to the error in the response. If err is not a hapiError
// then the default error message and default
func RespondError(w http.ResponseWriter, err error) error {
	return RespondErrorFallback(w, err, Config.DefaultStatusCode)
}

// RespondErrorFallback check if err is a type of hapiError. If it isn't, it will fallback
// to whatever status code you pass in.
func RespondErrorFallback(w http.ResponseWriter, err error, fallbackStatusCode int) error {
	statusCode := fallbackStatusCode
	message := Config.DefaultErrorMessage

	// check if err is hapi error
	hapiErr, ok := err.(hapiError)
	if ok {
		statusCode = hapiErr.GetStatusCode()
		message = hapiErr.GetMessage()
	}

	// if the message is still empty, get the default http status code message
	if message == "" {
		message = http.StatusText(statusCode)
	}

	errorResponse := errorResponse{
		Error: message,
	}

	if Config.ReturnRawError {
		errorResponse.RawError = err.Error()
	}

	return Respond(w, statusCode, errorResponse)
}

// RespondOK will marshal the payload and respond with a 200 status code.
func RespondOK(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusOK, payload)
}

// RespondBadRequest will marshal the error payload and respond with a 400 status code.
func RespondBadRequest(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusBadRequest, payload)
}

// RespondUnauthorized will marshal the error payload and respond with a 401 status code.
func RespondUnauthorized(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusUnauthorized, payload)
}

// RespondForbidden will marshal the error payload and respond with a 403 status code.
func RespondForbidden(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusForbidden, payload)
}

// RespondNotFound will marshal the error payload and respond with a 404 status code.
func RespondNotFound(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusNotFound, payload)
}

// RespondTeapot will marshal the error payload and respond with a 418 status code.
func RespondTeapot(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusTeapot, payload)
}

// RespondInternalError will marshal the error payload and respond with a 500 status code.
func RespondInternalError(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusInternalServerError, payload)
}
