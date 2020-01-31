package errors

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	testCases := []struct {
		desc     string
		hapiErr  hapiError
		expected string
	}{
		{
			desc: "Standard error",
			hapiErr: hapiError{
				err: errors.New("some standard error"),
			},
			expected: "some standard error",
		},
		{
			desc: "Wrapped error",
			hapiErr: hapiError{
				err: errors.Wrap(errors.New("some standard error"), "wrapped error"),
			},
			expected: "wrapped error: some standard error",
		},
		{
			desc: "Empty error",
			hapiErr: hapiError{
				err: errors.New(""),
			},
			expected: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.hapiErr.Error()

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSetMessageAndDetail(t *testing.T) {
	// This is needed because we need the same error message to be persisted in the tests
	daRealErr := New("this dumbass user is trying to access somewhere he shouldn't!!!")

	testCases := []struct {
		desc     string
		message  string
		err      error
		expected error
	}{
		{
			desc:    "Standard display message",
			message: "stephen is the GOAT",
			err:     hapiError{},
			expected: hapiError{
				message: "stephen is the GOAT",
			},
		},
		{
			desc:    "Make sure nothing else in hapiError changes",
			message: "You don't have permission to do that!",
			err: hapiError{
				errorType: Unauthorized,
				err:       daRealErr,
			},
			expected: hapiError{
				errorType: Unauthorized,
				err:       daRealErr,
				message:   "You don't have permission to do that!",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, ok := SetMessage(tc.err, tc.message)
			if !ok {
				assert.Fail(t, "SetDisplayMessage failed, error is not of type hapiError")
			}

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCastToHapiError(t *testing.T) {
	// Same reason as the previous test
	westCoastHapiError := hapiError{
		err: New("I heard you like playstation 2s..."),
	}

	testCases := []struct {
		desc              string
		err               error
		expectedHapiError hapiError
		expectedOk        bool
	}{
		{
			desc:              "Successful cast",
			err:               westCoastHapiError,
			expectedHapiError: westCoastHapiError,
			expectedOk:        true,
		},
		{
			desc:       "Unsuccessful cast",
			err:        errors.New("not our a hapiError, just the boring golang one"),
			expectedOk: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualHapiErr, actualOk := castToHapiError(tc.err)

			assert.Equal(t, tc.expectedOk, actualOk)

			if actualOk {
				assert.Equal(t, tc.expectedHapiError, actualHapiErr)
			}
		})
	}
}
