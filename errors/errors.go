package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Newf creates a NoType hapiError with formatted message
func Newf(format string, args ...interface{}) error {
	fmt.Println("creating newf error")

	return hapiError{
		errorType: NoType,
		err:       fmt.Errorf(format, args...),
	}
}

// New creates a NoType hapiError
func New(msg string) error {
	fmt.Println("creating new error")
	return hapiError{
		errorType: NoType,
		err:       Newf(msg),
	}
}

// Wrapf an error with a format string
func Wrapf(err error, format string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, format, args...)

	hapiErr, ok := castToHapiError(err)
	if !ok {
		return hapiError{
			errorType: NoType,
			err:       wrappedError,
		}
	}

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
