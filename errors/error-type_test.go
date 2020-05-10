package errors

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func compareErrors(t *testing.T, expectedErr error, actualErr error) {
	expected := castToHapiError(expectedErr)
	actual := castToHapiError(actualErr)

	assert.Equal(t, expected.Error(), actual.Error())
	assert.Equal(t, expected.ErrorType, actual.ErrorType)
	assert.Equal(t, expected.Message, actual.Message)
}

func TestErrorTypeNewf(t *testing.T) {
	testCases := []struct {
		desc      string
		errorType ErrorType
		msg       string
		args      []interface{}
		expected  error
	}{
		{
			desc:      "Standard error",
			errorType: ImATeapot,
			msg:       "Mindhunter is a pretty good show",
			expected: HapiError{
				ErrorType: ImATeapot,
				Err:       errors.New("Mindhunter is a pretty good show"),
			},
		},
		{
			desc:      "Error with format",
			errorType: InternalServerError,
			msg:       "%d is the answer to life",
			args:      []interface{}{42},
			expected: HapiError{
				ErrorType: InternalServerError,
				Err:       errors.Errorf("%d is the answer to life", 42),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.errorType.Newf(tc.msg, tc.args...)

			compareErrors(t, tc.expected, actual)
		})
	}
}

func TestErrorTypeNew(t *testing.T) {
	testCases := []struct {
		desc      string
		errorType ErrorType
		msg       string
		expected  error
	}{
		{
			desc:      "Standard error",
			errorType: Forbidden,
			msg:       "YOU SHALL NOT PASS",
			expected: HapiError{
				ErrorType: Forbidden,
				Err:       New("YOU SHALL NOT PASS"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.errorType.New(tc.msg)

			compareErrors(t, tc.expected, actual)
		})
	}
}

func TestErrorTypeWrapf(t *testing.T) {
	testCases := []struct {
		desc      string
		errorType ErrorType
		errToWrap error
		msg       string
		args      []interface{}
		expected  error
	}{
		{
			desc:      "Standard wrapped error",
			errorType: ImATeapot,
			errToWrap: errors.New("dis gone get wrapped"),
			msg:       "dis da wrapper",
			expected: HapiError{
				ErrorType: ImATeapot,
				Err:       errors.Wrap(errors.New("dis gone get wrapped"), "dis da wrapper"),
			},
		},
		{
			desc:      "Error with format",
			errorType: InternalServerError,
			errToWrap: errors.New("42 is the answer to life"),
			msg:       "%d is still the answer to life",
			args:      []interface{}{42},
			expected: HapiError{
				ErrorType: InternalServerError,
				Err:       errors.Wrapf(errors.New("42 is the answer to life"), "%d is still the answer to life", 42),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.errorType.Wrapf(tc.errToWrap, tc.msg, tc.args...)

			assert.Equal(t, tc.expected.Error(), actual.Error())
		})
	}
}

func TestErrorTypeWrap(t *testing.T) {
	testCases := []struct {
		desc      string
		errType   ErrorType
		errToWrap error
		msg       string
		expected  error
	}{
		{
			desc:      "Standard wrapped error",
			errType:   NotFound,
			errToWrap: errors.New("og error"),
			msg:       "gonna wrap dat og error",
			expected: HapiError{
				ErrorType: NotFound,
				Err:       errors.Wrap(errors.New("og error"), "gonna wrap dat og error"),
			},
		},
		{
			desc:    "Wrap error with different error type",
			errType: NotFound,
			errToWrap: HapiError{
				ErrorType: Unauthorized,
				Err:       errors.New("og error"),
			},
			msg: "gonna wrap dat og error",
			expected: HapiError{
				ErrorType: NotFound,
				Err:       errors.Wrap(errors.New("og error"), "gonna wrap dat og error"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.errType.Wrap(tc.errToWrap, tc.msg)

			compareErrors(t, tc.expected, actual)
		})
	}
}

func TestHapiErrorCast(t *testing.T) {
	testCases := []struct {
		desc      string
		errToCast error
		expected  error
	}{
		{
			desc:      "Standard library error",
			errToCast: errors.New("og error"),
			expected: HapiError{
				ErrorType: NoType,
				Err:       errors.New("og error"),
			},
		},
		{
			desc:      "Hapi error cast",
			errToCast: ImATeapot.New("some awesome hapi error"),
			expected: HapiError{
				ErrorType: ImATeapot,
				Err:       errors.New("some awesome hapi error"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := Cast(tc.errToCast)

			compareErrors(t, tc.expected, actual)
		})
	}
}
