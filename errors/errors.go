package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Newf creates a NoType hapiError with formatted message
func Newf(format string, args ...interface{}) error {
	return hapiError{
		errorType: NoType,
		err:       fmt.Errorf(format, args...),
	}
}

// New creates a NoType hapiError
func New(msg string) error {
	return hapiError{
		errorType: NoType,
		err:       errors.New(msg),
	}
}

// Wrapf an error with a format string
func Wrapf(err error, format string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, format, args...)

	hapiErr := castToHapiError(err)

	return hapiError{
		errorType: hapiErr.errorType,
		err:       wrappedError,
		message:   hapiErr.message,
	}
}

// Wrap an error with a string
func Wrap(err error, message string) error {
	return Wrapf(err, message)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}
