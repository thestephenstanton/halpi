package hapi

import (
	"encoding/json"
	"net/http"

	"github.com/thestephenstanton/hapi/errors"
)

// Respond will return a json marshaled to the client. Error will only return
// if payload was not able to be marshaled
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

// 200
func RespondOK(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusOK, payload)
}

// 400
func RespondBadRequest(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusBadRequest, payload)
}

// 401
func RespondUnauthorized(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusUnauthorized, payload)
}

// 403
func RespondForbidden(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusForbidden, payload)
}

// 404
func RespondNotFound(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusNotFound, payload)
}

// 418
func RespondTeapot(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusTeapot, payload)
}

// 500
func RespondInternalError(w http.ResponseWriter, payload interface{}) error {
	return Respond(w, http.StatusInternalServerError, payload)
}
