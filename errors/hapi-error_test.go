package errors

import (
	"testing"

	goerrors "github.com/pkg/errors"
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
				err: goerrors.New("some standard error"),
			},
			expected: "some standard error",
		},
		{
			desc: "Wrapped error",
			hapiErr: hapiError{
				err: goerrors.Wrap(goerrors.New("some standard error"), "wrapped error"),
			},
			expected: "wrapped error: some standard error",
		},
		{
			desc: "Empty error",
			hapiErr: hapiError{
				err: goerrors.New(""),
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

func TestSetMessage(t *testing.T) {
	// This is needed because we need the same error message to be persisted in the tests
	daRealErr := New("this dumbass user is trying to access somewhere he shouldn't!!!")

	testCases := []struct {
		desc     string
		message  string
		err      error
		expected string
	}{
		{
			desc:     "Standard display message",
			message:  "stephen is the GOAT",
			err:      hapiError{},
			expected: "stephen is the GOAT",
		},
		{
			desc:    "Make sure nothing else in hapiError changes",
			message: "You don't have permission to do that!",
			err: hapiError{
				errorType: Unauthorized,
				err:       daRealErr,
			},
			expected: "You don't have permission to do that!",
		},
		{
			desc:     "Get regular error and return back a hapi error with message set",
			message:  "You don't have permission to do that!",
			err:      goerrors.New("some regular error"),
			expected: "You don't have permission to do that!",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := SetMessage(tc.err, tc.message)

			actualErr := castToHapiError(err)

			assert.Equal(t, tc.expected, actualErr.GetMessage())
		})
	}
}

func TestCastToHapiError(t *testing.T) {
	// Same reason as the previous test
	westCoastHapiError := hapiError{
		err: New("I heard you like playstation 2s..."),
	}

	standardLibraryErr := goerrors.New("not our a hapiError, just the boring golang one, but we will fix")

	testCases := []struct {
		desc              string
		err               error
		expectedHapiError hapiError
	}{
		{
			desc:              "Hapi error, happy cast",
			err:               westCoastHapiError,
			expectedHapiError: westCoastHapiError,
		},
		{
			desc: "Not hapi error, cast to a hapi error",
			err:  standardLibraryErr,
			expectedHapiError: hapiError{
				errorType: NoType,
				err:       standardLibraryErr,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualHapiErr := castToHapiError(tc.err)

			assert.Equal(t, tc.expectedHapiError, actualHapiErr)
		})
	}
}
