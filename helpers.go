package hapi

import (
	"encoding/json"
	"net/http"

	"github.com/thestephenstanton/hapi/errors"
)

// GetQueryParam will return the first value for the key provided and true. If no key was found OR
// if the key's value was empty, it will return an empty string and false
func GetQueryParam(request *http.Request, key string) (string, bool) {
	value := request.URL.Query().Get(key)
	if value == "" {
		return "", false
	}

	return value, true
}

// UnmarshalBody will unmarshal the request's body into the interface provided
func UnmarshalBody(request *http.Request, v interface{}, opts ...UnmarshalOption) error {
	for _, opt := range opts {
		opt(request)
	}

	err := json.NewDecoder(request.Body).Decode(&v)
	if err != nil {
		if len(opts) > 0 && err.Error() == requestBodyTooLargeError {
			return errors.TooLarge.Wrap(err, "request body is too large")
		}

		return errors.BadRequest.Wrap(err, "request body is not proper json")
	}

	return nil
}

// UnmarshalOption is an option with given the request for unmarshalling
type UnmarshalOption func(r *http.Request)

const requestBodyTooLargeError = "http: request body too large"

// WithMaxSize will wrap the request in http.MaxBytesReader before trying to unmarshal
func WithMaxSize(w http.ResponseWriter, maxBytes int64) UnmarshalOption {
	return func(r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	}
}
