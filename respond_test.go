package hapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	goerrors "github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/thestephenstanton/hapi/errors"
)

// This helps out with the two basic responders: Respond and RespondError.
func respondHandler(statusCode int, payload interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := Respond(w, statusCode, payload)
		if err != nil {
			panic(err.Error())
		}
	}
}

func respondErrorHandler(err error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := RespondError(w, err)
		if err != nil {
			panic(err.Error())
		}
	}
}

// This helps with all the helper responders like RespondOK, RespondNotFound, etc...
func helperResponderHandlerHelper(helperResponder func(http.ResponseWriter, interface{}) error, payload interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helperResponder(w, payload)
		if err != nil {
			panic(err.Error())
		}
	}
}

func toPointer(s string) *string {
	return &s
}

func TestRespond(t *testing.T) {
	testCases := []struct {
		desc               string
		statusCode         int
		payload            interface{}
		expectedStatusCode int
		expectedBody       string
	}{
		{
			desc:               "basic response",
			statusCode:         42,
			payload:            "hello world",
			expectedStatusCode: 42,
			expectedBody:       `"hello world"`,
		},
		{
			desc:               "nil payload",
			statusCode:         741,
			payload:            nil,
			expectedStatusCode: 741,
			expectedBody:       `null`,
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
			expectedBody:       `{"foo":"bar"}`,
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
			actualBody := recorder.Body.String()

			assert.Equal(t, tc.expectedStatusCode, actualStatusCode)
			assert.Equal(t, tc.expectedBody, actualBody)
		})
	}
}

func TestRespondError(t *testing.T) {
	testCases := []struct {
		desc               string
		err                error
		expectedStatusCode int
		expectedBody       string
	}{
		{
			desc:               "respond with hapi error",
			err:                errors.Forbidden.New("some message you want your client to see"),
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"error":"some message you want your client to see"}`,
		},
		{
			desc: "custom set message",
			err: errors.SetMessage(
				errors.Forbidden.New("info you might not want client to see"),
				"some message you want your client to see",
			),
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"error":"some message you want your client to see"}`,
		},
		{
			desc:               "standard go error",
			err:                goerrors.New("some go error"),
			expectedStatusCode: Config.DefaultStatusCode,
			expectedBody:       fmt.Sprintf(`{"error":"%s"}`, Config.DefaultErrorMessage),
		},
		{
			desc: "custom error matching hapiError interface",
			err: customHapiError{
				statusCode: 666,
				message:    "custom error someone created",
			},
			expectedStatusCode: 666,
			expectedBody:       `{"error":"custom error someone created"}`,
		},
		{
			desc:               "nil error",
			err:                nil,
			expectedStatusCode: Config.DefaultStatusCode,
			expectedBody:       fmt.Sprintf(`{"error":"%s"}`, Config.DefaultErrorMessage),
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
			assert.Equal(t, tc.expectedBody, actualBody)
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

func TestChangeDefaults(t *testing.T) {
	originalConfig := Config

	testCases := []struct {
		desc                   string
		newDefaultStatusCode   int
		newDefaultErrorMessage string
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
