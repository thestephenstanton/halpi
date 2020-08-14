package hapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRequestWithURL(t *testing.T, url string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Error("failed to create new request", err)
		t.FailNow()
	}

	return request
}

func TestGetQueryParam(t *testing.T) {
	testCases := []struct {
		desc           string
		request        *http.Request
		key            string
		expectedResult string
		expectedFound  bool
	}{
		{
			desc:           "standard query param",
			request:        newRequestWithURL(t, "http://test.com/foo?bar=fubar"),
			key:            "bar",
			expectedResult: "fubar",
			expectedFound:  true,
		},
		{
			desc:          "no key found",
			request:       newRequestWithURL(t, "http://test.com/foo?bar=fubar"),
			key:           "nokey",
			expectedFound: false,
		},
		{
			desc:           "find first key",
			request:        newRequestWithURL(t, "http://test.com/foo?bar=first&bar=second"),
			key:            "bar",
			expectedResult: "first",
			expectedFound:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result, found := GetQueryParam(tc.request, tc.key)

			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedFound, found)
		})
	}
}

type testStruct struct {
	Text string
}

func TestUnmarshalBody(t *testing.T) {
	testCases := []struct {
		desc           string
		request        *http.Request
		opts           []UnmarshalOption
		expectedStruct testStruct
		shouldError    bool
		expectedError  string
	}{
		{
			desc: "normal unmarshal",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader(json.RawMessage(`{"text":"hello world"}`))),
			},
			expectedStruct: testStruct{
				Text: "hello world",
			},
		},
		{
			desc: "failed unmarshal",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader([]byte("bad json"))),
			},
			shouldError:   true,
			expectedError: "request body is not proper json: invalid character 'b' looking for beginning of value",
		},
		{
			desc: "max bytes option success",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader(json.RawMessage(`{"text":"hello world"}`))),
			},
			opts: []UnmarshalOption{
				WithMaxSize(nil, 100),
			},
			expectedStruct: testStruct{
				Text: "hello world",
			},
		},
		{
			desc: "max bytes option failure",
			request: &http.Request{
				Body: ioutil.NopCloser(bytes.NewReader(json.RawMessage(`{"text":"hello world"}`))),
			},
			opts: []UnmarshalOption{
				WithMaxSize(nil, 1),
			},
			shouldError:   true,
			expectedError: "request body is too large: http: request body too large",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var actualStruct testStruct
			err := UnmarshalBody(tc.request, &actualStruct, tc.opts...)
			if tc.shouldError && err == nil {
				assert.Fail(t, "should have errored but didn't")
			}
			if err != nil {
				assert.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedStruct, actualStruct)
		})
	}
}
