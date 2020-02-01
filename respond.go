package hapi

import (
	"encoding/json"
	"net/http"

	"github.com/thestephenstanton/hapi/errors"
)

var DefaultNotFoundErr = errors.BadRequest.New("this is not the endpoint you are looking for...")

func RespondOK(w http.ResponseWriter, payload interface{}) {
	Respond(w, NewResponseEnvelope(http.StatusOK, payload))
}

func RespondNotFound(w http.ResponseWriter) {
	Respond(w, NewErrorEnvelope(http.StatusNotFound, DefaultNotFoundErr))
}

func RespondBadRequest(w http.ResponseWriter, err error) {
	Respond(w, NewErrorEnvelope(http.StatusBadRequest, err))
}

func RespondInternalError(w http.ResponseWriter, err error) {
	Respond(w, NewErrorEnvelope(http.StatusInternalServerError, err))
}

func Respond(w http.ResponseWriter, envelope ResponseEnvelope) {
	bytes, err := json.Marshal(envelope)
	if err != nil {
		err = errors.New("failed to marshal json, unhapi")
		RespondInternalError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(envelope.StatusCode)

	// Prevents "null" from being sent back
	if string(bytes) != "null" {
		w.Write(bytes)
	}
}
