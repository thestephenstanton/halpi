package halpi

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func ReplyOK(writer http.ResponseWriter, payload interface{}) {
	Reply(writer, http.StatusOK, payload)
}

func ReplyBadRequest(writer http.ResponseWriter, payload interface{}) {
	Reply(writer, http.StatusBadRequest, payload)
}

func ReplyError(writer http.ResponseWriter, err error) {
	Reply(writer, http.StatusInternalServerError, err.Error())
}

func Reply(writer http.ResponseWriter, statusCode int, payload interface{}) {
	json, err := json.Marshal(payload)
	if err != nil {
		err = errors.Wrap(err, "unable to marshal json")

		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application-json")
	writer.WriteHeader(statusCode)

	// Prevents "null" from being sent back
	if string(json) != "null" {
		writer.Write(json)
	}
}
