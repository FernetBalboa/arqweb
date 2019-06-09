package apierror

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

const (
	NoType = ErrorType(iota)
	BadRequest
	Forbidden
	NotFound
	Internal
	External
	Parsing
)

type ErrorType uint

// wrappingError represents an error that may have type.
// It also holds another error as it's cause.
type wrappingError struct {
	errorType     ErrorType
	originalError error
}

// Error returns the message of a wrappingError
func (ce wrappingError) Error() string {
	return ce.originalError.Error()
}

// Cause returns the cause of the custom error, recursively.
func (ce wrappingError) Cause() error {
	return errors.Cause(ce.originalError)
}

func (ce wrappingError) Format(f fmt.State, c rune) {
	original := ce.originalError
	if customErr, ok := original.(fmt.Formatter); ok {
		customErr.Format(f, c)
	} else {
		_, _ = io.WriteString(f, ce.Error())
	}
}

// New creates a new wrappingError
func (et ErrorType) New(msg string) error {
	return wrappingError{errorType: et, originalError: errors.New(msg)}
}

// New creates a new wrappingError with formatted message
func (et ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)

	return wrappingError{errorType: et, originalError: err}
}

// Wrap creates a new wrapped error
func (et ErrorType) Wrap(err error, msg string) error {
	return et.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (et ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return wrappingError{errorType: et, originalError: newErr}
}

// New creates a no type error
func New(msg string) error {
	return wrappingError{errorType: NoType, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return wrappingError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap wraps an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(wrappingError); ok {
		return wrappingError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
		}
	}

	return wrappingError{errorType: NoType, originalError: wrappedError}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(wrappingError); ok {
		return customErr.errorType
	}

	return NoType
}
