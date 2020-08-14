package hapi

import "net/http"

// Config configs hapi
var Config = struct {
	DefaultErrorMessage string
	DefaultStatusCode   int
	ReturnNulls         bool
	ReturnRawError      bool
}{
	DefaultErrorMessage: "uh oh, something went wrong, please try again later",
	DefaultStatusCode:   http.StatusInternalServerError,
	ReturnNulls:         false,
	ReturnRawError:      false,
}
