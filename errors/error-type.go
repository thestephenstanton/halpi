package errors

import (
	"net/http"

	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// NoType error
	NoType ErrorType = iota + 1

	// BadRequest 400 error
	BadRequest

	// Unauthorized 401 error
	Unauthorized

	// Forbidden 403 error
	Forbidden

	// NotFound 404 error
	NotFound

	// ImATeapot 418
	ImATeapot

	// InternalServerError 500 error
	InternalServerError
)

// Newf creates a new hapiError with formatted message
func (errorType ErrorType) Newf(message string, args ...interface{}) error {
	return hapiError{
		errorType: errorType,
		err:       errors.Errorf(message, args...),
	}
}

// New creates a new hapiError
func (errorType ErrorType) New(message string) error {
	return hapiError{
		errorType: errorType,
		err:       errorType.Newf(message),
	}
}

// Wrapf creates a new wrapped hapiError with formatted message
func (errorType ErrorType) Wrapf(err error, message string, args ...interface{}) error {
	return hapiError{
		errorType: errorType,
		err:       errors.Wrapf(err, message, args...),
	}
}

// Wrap creates a new wrapped hapiError
func (errorType ErrorType) Wrap(err error, message string) error {
	return errorType.Wrapf(err, message)
}

func getStatusCode(errorType ErrorType) int {
	switch errorType {
	case BadRequest:
		return http.StatusBadRequest // 400
	case Unauthorized:
		return http.StatusUnauthorized // 401
	case Forbidden:
		return http.StatusForbidden // 403
	case NotFound:
		return http.StatusNotFound // 404
	case ImATeapot:
		return http.StatusTeapot // 418
	case InternalServerError:
		return http.StatusInternalServerError // 500
	default:
		return http.StatusInternalServerError // 500
	}
}
