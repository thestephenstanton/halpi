package halpi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func GetQueryParam(request *http.Request, key string) (string, error) {
	vars := mux.Vars(request)
	stringKey, ok := vars[key]
	if !ok {
		return "", fmt.Errorf("no '%s' found", key)
	}

	return stringKey, nil
}

func UnmarshalBody(request *http.Request, model interface{}) error {
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&model)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal")
	}

	return nil
}
