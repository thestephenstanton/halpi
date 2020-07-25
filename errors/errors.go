package errors

import (
	"github.com/pkg/errors"
)

// Newf creates a NoType hapiError with formatted message
func Newf(format string, args ...interface{}) error {
	return NoType.Newf(format, args...)
}

// New creates a NoType hapiError
func New(message string) error {
	return Newf(message)
}

// Wrapf an error with a format string
func Wrapf(err error, format string, args ...interface{}) error {
	return NoType.Wrapf(err, format, args...)
}

// Wrap an error with a string
func Wrap(err error, message string) error {
	return Wrapf(err, message)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}
