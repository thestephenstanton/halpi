package hapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReponseEnvelope(t *testing.T) {
	testCases := []struct {
		desc       string
		statusCode int
		data       interface{}
		expected   ResponseEnvelope
	}{
		{
			desc:       "Normal response envelope",
			statusCode: 42,
			data:       "the world is on fire",
			expected: ResponseEnvelope{
				StatusCode: 42,
				Data:       "the world is on fire",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := NewResponseEnvelope(tc.statusCode, tc.data)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestNewErrorEnvelope(t *testing.T) {
	testCases := []struct {
		desc       string
		statusCode int
		err        error
		expected   ResponseEnvelope
	}{
		{
			desc:       "Normal error envelope with hapi error",
			statusCode: 69,
			// err:        ,
			expected: ResponseEnvelope{
				StatusCode: 69,
				Error: ErrorEnvelope{
					Message: "you done fucked up a a ron",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// actual := NewErrorEnvelope(tc.statusCode, tc.data)
			// assert.Equal(t, tc.expected, actual)
		})
	}
}
