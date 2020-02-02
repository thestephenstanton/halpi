package hapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	goerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thestephenstanton/hapi/errors"
)

// This helps out with the two basic responders: Respond and RespondError.
func basicResponderHandlerHelper(basicResponder func(http.ResponseWriter, int, interface{}) error, statusCode int, payload interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := basicResponder(w, statusCode, payload)
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

func TestRespond(t *testing.T) {
	testCases := []struct {
		desc                  string
		configUseHapiEnvelope bool
		statusCode            int
		payload               interface{}
		expectedStatusCode    int
		expectedBody          string
	}{
		{
			desc:                  "Basic response, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            42,
			payload:               "hello world",
			expectedStatusCode:    42,
			expectedBody:          `{"statusCode":42,"data":"hello world"}`,
		},
		{
			desc:                  "Basic response, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            42,
			payload:               "hello world",
			expectedStatusCode:    42,
			expectedBody:          `"hello world"`,
		},
		{
			desc:                  "Nil payload, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            741,
			payload:               nil,
			expectedStatusCode:    741,
			expectedBody:          `{"statusCode":741}`,
		},
		{
			desc:                  "Nil payload, with hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            741,
			payload:               nil,
			expectedStatusCode:    741,
			expectedBody:          ``,
		},
		{
			desc:                  "Custom payload, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            951,
			payload: struct {
				Foo string `json:"foo"`
			}{
				Foo: "bar",
			},
			expectedStatusCode: 951,
			expectedBody:       `{"statusCode":951,"data":{"foo":"bar"}}`,
		},
		{
			desc:                  "Custom payload, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            951,
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
			Config.UseHapiEnvelopes = tc.configUseHapiEnvelope

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(basicResponderHandlerHelper(Respond, tc.statusCode, tc.payload))

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
		desc                  string
		configUseHapiEnvelope bool
		statusCode            int
		payload               interface{}
		expectedStatusCode    int
		expectedBody          string
	}{
		{
			desc:                  "hapi error, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            123, // overwritten by hapi error
			payload: errors.SetMessage(
				errors.Forbidden.New("info you might not want client to see"),
				"some message you want your client to see",
			),
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `{"statusCode":403,"error":{"message":"some message you want your client to see"}}`,
		},
		{
			desc:                  "hapi error, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            123, // overwritten by hapi error
			payload: errors.SetMessage(
				errors.Forbidden.New("info you might not want client to see"),
				"some message you want your client to see",
			),
			expectedStatusCode: http.StatusForbidden,
			expectedBody:       `"some message you want your client to see"`,
		},
		{
			desc:                  "standard go error, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            420,
			payload:               goerrors.New("some new error created that you want your client to see"),
			expectedStatusCode:    420,
			expectedBody:          `{"statusCode":420,"error":{"message":"some new error created that you want your client to see"}}`,
		},
		{
			desc:                  "standard go error, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            666,
			payload:               goerrors.New("some new error created that you want your client to see"),
			expectedStatusCode:    666,
			expectedBody:          `"some new error created that you want your client to see"`,
		},
		{
			desc:                  "custom hapi error, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            432, // will be overwriten by custom hapi error
			payload: customHapiError{
				statusCode: 789,
				message:    "some custom message",
			},
			expectedStatusCode: 789,
			expectedBody:       `{"statusCode":789,"error":{"message":"some custom message"}}`,
		},
		{
			desc:                  "custom hapi error, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            432, // will be overwriten by custom hapi error
			payload: customHapiError{
				statusCode: 789,
				message:    "some custom message",
			},
			expectedStatusCode: 789,
			expectedBody:       `"some custom message"`,
		},
		{
			desc:                  "no error, just a payload, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            639,
			payload:               "some error message",
			expectedStatusCode:    639,
			expectedBody:          `{"statusCode":639,"error":"some error message"}`,
		},
		{
			desc:                  "no error, just a payload, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            298,
			payload:               "some error message",
			expectedStatusCode:    298,
			expectedBody:          `"some error message"`,
		},
		{
			desc:                  "nil payload, with hapi envelope",
			configUseHapiEnvelope: true,
			statusCode:            489,
			payload:               nil,
			expectedStatusCode:    489,
			expectedBody:          `{"statusCode":489}`,
		},
		{
			desc:                  "nil payload, no hapi envelope",
			configUseHapiEnvelope: false,
			statusCode:            165,
			payload:               nil,
			expectedStatusCode:    165,
			expectedBody:          ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			Config.UseHapiEnvelopes = tc.configUseHapiEnvelope

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(basicResponderHandlerHelper(RespondError, tc.statusCode, tc.payload))

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
	Config.UseHapiEnvelopes = false
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
