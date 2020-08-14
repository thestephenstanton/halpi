package hapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	goerrors "github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/thestephenstanton/hapi/errors"
)

// helps with TestRespond
func respondHandler(statusCode int, payload interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := Respond(w, statusCode, payload)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestRespond(t *testing.T) {
	testCases := []struct {
		desc               string
		statusCode         int
		payload            interface{}
		expectedStatusCode int
		expectedBody       []byte
	}{
		{
			desc:               "basic response",
			statusCode:         42,
			payload:            "hello world",
			expectedStatusCode: 42,
			expectedBody:       []byte(`"hello world"`),
		},
		{
			desc:               "nil payload",
			statusCode:         741,
			payload:            nil,
			expectedStatusCode: 741,
			expectedBody:       nil,
		},
		{
			desc:       "custom payload",
			statusCode: 951,
			payload: struct {
				Foo string `json:"foo"`
			}{
				Foo: "bar",
			},
			expectedStatusCode: 951,
			expectedBody:       json.RawMessage(`{"foo":"bar"}`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(respondHandler(tc.statusCode, tc.payload))

			// We don't care about the request, we just care about the response.
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			actualStatusCode := recorder.Code
			actualBody := recorder.Body.Bytes()

			assert.Equal(t, tc.expectedStatusCode, actualStatusCode)
			assert.Equal(t, tc.expectedBody, actualBody)
		})
	}
}

// heals with TestRespondError
func respondErrorHandler(err error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := RespondError(w, err)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestRespondError(t *testing.T) {
	testCases := []struct {
		desc               string
		err                error
		expectedStatusCode int
		expectedBody       json.RawMessage
	}{
		{
			desc:               "respond with hapi error",
			err:                errors.Forbidden.New("some message you want your client to see"),
			expectedStatusCode: http.StatusForbidden,
			expectedBody: json.RawMessage(`
			{
				"error": "some message you want your client to see"
			}
			`),
		},
		{
			desc: "custom set message",
			err: errors.SetMessage(
				errors.Forbidden.New("info you might not want client to see"),
				"some message you want your client to see",
			),
			expectedStatusCode: http.StatusForbidden,
			expectedBody: json.RawMessage(`
			{
				"error": "some message you want your client to see"
			}
			`),
		},
		{
			desc:               "another way to set message",
			err:                errors.Forbidden.New("info you might not want client to see").SetMessage("some message you want your client to see"),
			expectedStatusCode: http.StatusForbidden,
			expectedBody: json.RawMessage(`
			{
				"error": "some message you want your client to see"
			}
			`),
		},
		{
			desc:               "standard go error",
			err:                goerrors.New("some go error"),
			expectedStatusCode: Config.DefaultStatusCode,
			expectedBody:       json.RawMessage(fmt.Sprintf(`{"error":"%s"}`, Config.DefaultErrorMessage)),
		},
		{
			desc:               "set message on standard go errors",
			err:                errors.SetMessage(goerrors.New("some go error"), "created by errors.SetMessage()"),
			expectedStatusCode: Config.DefaultStatusCode,
			expectedBody: json.RawMessage(`
			{
				"error": "created by errors.SetMessage()"
			}
			`),
		},
		{
			desc: "custom error matching hapiError interface",
			err: customHapiError{
				statusCode: 666,
				message:    "custom error someone created",
			},
			expectedStatusCode: 666,
			expectedBody: json.RawMessage(`
			{
				"error": "custom error someone created"
			}
			`),
		},
		{
			desc:               "nil error",
			err:                nil,
			expectedStatusCode: Config.DefaultStatusCode,
			expectedBody:       json.RawMessage(fmt.Sprintf(`{"error":"%s"}`, Config.DefaultErrorMessage)),
		},
		{
			desc:               "hapi error wrapped with hapi error",
			err:                errors.Unauthorized.Wrap(errors.ImATeapot.New("initial error"), "wrapping error"),
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody: json.RawMessage(`
			{
				"error": "wrapping error"
			}
			`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(respondErrorHandler(tc.err))

			// We don't care about the request, we just care about the response.
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			actualStatusCode := recorder.Code
			actualBody := recorder.Body.String()

			assert.Equal(t, tc.expectedStatusCode, actualStatusCode)
			assert.JSONEq(t, string(tc.expectedBody), string(actualBody))
		})
	}
}

type customHapiError struct {
	statusCode int
	message    string
}

func (c customHapiError) Error() string {
	return "something went wrong"
}

func (c customHapiError) GetStatusCode() int {
	return c.statusCode
}

func (c customHapiError) GetMessage() string {
	return c.message
}

// helps with TestRespondHelpers
func helperResponderHandlerHelper(helperResponder func(http.ResponseWriter, interface{}) error, payload interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helperResponder(w, payload)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestRespondHelpers(t *testing.T) {
	// We are just going to assume these are always the same since we test this thoroughly
	// in TestRespond and TestRespondError
	payload := "hello world"
	expectedBody := `"hello world"`

	testCases := []struct {
		desc               string
		respond            func(http.ResponseWriter, interface{}) error
		expectedStatusCode int
	}{
		{
			desc:               "Test OK",
			respond:            RespondOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			desc:               "Test Bad Request",
			respond:            RespondBadRequest,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			desc:               "Test Unauthorized",
			respond:            RespondUnauthorized,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			desc:               "Test Forbidden",
			respond:            RespondForbidden,
			expectedStatusCode: http.StatusForbidden,
		},
		{
			desc:               "Test NotFound",
			respond:            RespondNotFound,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			desc:               "Test Teapot",
			respond:            RespondTeapot,
			expectedStatusCode: http.StatusTeapot,
		},
		{
			desc:               "Test TooLarge",
			respond:            RespondTooLarge,
			expectedStatusCode: http.StatusRequestEntityTooLarge,
		},
		{
			desc:               "Test InternalError",
			respond:            RespondInternalError,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(helperResponderHandlerHelper(tc.respond, payload))

			// We don't care about the request, we just care about the response.
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			actualStatusCode := recorder.Code
			actualBody := recorder.Body.String()

			assert.Equal(t, tc.expectedStatusCode, actualStatusCode)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}

func TestChangeConfigDefaults(t *testing.T) {
	originalConfig := Config

	testCases := []struct {
		desc                   string
		newDefaultStatusCode   int
		newDefaultErrorMessage string
		newReturnNulls         bool
		newReturnRawError      bool
		expectedStatusCode     int
		expectedBody           string
	}{
		{
			desc:                   "change default status code",
			newDefaultStatusCode:   123,
			newDefaultErrorMessage: originalConfig.DefaultErrorMessage,
			expectedStatusCode:     123,
			expectedBody:           fmt.Sprintf(`{"error":"%s"}`, Config.DefaultErrorMessage),
		},
		{
			desc:                   "change default message",
			newDefaultStatusCode:   originalConfig.DefaultStatusCode,
			newDefaultErrorMessage: "new default message",
			expectedStatusCode:     originalConfig.DefaultStatusCode,
			expectedBody:           fmt.Sprintf(`{"error":"%s"}`, "new default message"),
		},
		{
			desc:                   "empty default error message",
			newDefaultStatusCode:   http.StatusTeapot,
			newDefaultErrorMessage: "",
			expectedStatusCode:     http.StatusTeapot,
			expectedBody:           fmt.Sprintf(`{"error":"%s"}`, "I'm a teapot"),
		},
		{
			desc:                   "empty default error message with not real status code",
			newDefaultStatusCode:   42,
			newDefaultErrorMessage: "",
			expectedStatusCode:     42,
			expectedBody:           fmt.Sprintf(`{"error":"%s"}`, ""),
		},
		{
			desc:                   "return raw error",
			newDefaultStatusCode:   originalConfig.DefaultStatusCode,
			newDefaultErrorMessage: originalConfig.DefaultErrorMessage,
			newReturnRawError:      true,
			expectedStatusCode:     originalConfig.DefaultStatusCode,
			expectedBody:           fmt.Sprintf(`{"error":"%s","rawError":"%s"}`, Config.DefaultErrorMessage, "detailed error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			Config.DefaultStatusCode = tc.newDefaultStatusCode
			Config.DefaultErrorMessage = tc.newDefaultErrorMessage
			Config.ReturnRawError = tc.newReturnRawError

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(respondErrorHandler(goerrors.New("detailed error")))

			// We don't care about the request, we just care about the response.
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			actualStatusCode := recorder.Code
			actualBody := recorder.Body.String()

			assert.Equal(t, tc.expectedStatusCode, actualStatusCode)
			assert.Equal(t, tc.expectedBody, actualBody)

			// reset config
			Config = originalConfig
		})
	}
}
