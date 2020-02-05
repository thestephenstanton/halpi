package hapi

import (
	"encoding/json"
	"net/http"

	"github.com/thestephenstanton/hapi/errors"
)

// Respond will marshal and return the payload to the client with a given status code.
func Respond(w http.ResponseWriter, statusCode int, payload interface{}) error {
	if Config.UseHapiEnvelopes {
		payload = NewResponseEnvelope(statusCode, payload)
	}

	return respond(w, statusCode, payload)
}

func respond(w http.ResponseWriter, statusCode int, payload interface{}) error {
	var bytes []byte
	var err error

	if payload != nil {
		bytes, err = json.Marshal(payload)
		if err != nil {
			return errors.InternalServerError.Wrap(err, "failed to marshal payload")
		}
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(statusCode)
	w.Write(bytes)

	return nil
}

// RespondError will marshal and return the error payload to the client with a given status code.
func RespondError(w http.ResponseWriter, statusCode int, payload interface{}) error {
	// Check if it is a hapi error
	hapiErr, ok := payload.(hapiError)
	if ok {
		if Config.UseHapiEnvelopes {
			payload = NewErrorEnvelope(statusCode, hapiErr.(error))
		} else {
			payload = hapiErr.GetMessage()
		}

		return respond(w, hapiErr.GetStatusCode(), payload)
	}

	// Check if just normal error
	regularErr, ok := payload.(error)
	if ok {
		if Config.UseHapiEnvelopes {
			payload = NewErrorEnvelope(statusCode, regularErr)
		} else {
			payload = regularErr.Error()
		}

		return respond(w, statusCode, payload)
	}

	// Otherwise, send the payload
	if Config.UseHapiEnvelopes {
		payload = NewErrorEnvelope(statusCode, payload)
	}

	return respond(w, statusCode, payload)
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
