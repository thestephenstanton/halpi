package hapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thestephenstanton/hapi/errors"
)

func GetQueryParam(request *http.Request, key string) (string, error) {
	vars := mux.Vars(request)
	stringKey, ok := vars[key]
	if !ok {
		return "", errors.InternalServerError.Newf("query param '%s' not found, unhapi", key)
	}

	return stringKey, nil
}

func UnmarshalBody(request *http.Request, v interface{}) error {
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&v)
	if err != nil {
		return errors.InternalServerError.Wrap(err, "failed to unmarshal, unhapi")
	}

	return nil
}
