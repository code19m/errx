package errx

import (
	"errors"
)

// Type defines the different categories of errors that can be represented by an ErrorX.
type Type int8

// M is a shorthand for a map of string key-value pairs.
type M map[string]string

const (
	T_Internal   Type = iota // Internal errors indicate unexpected issues within the application.
	T_Validation             // Validation errors occur when user input does not meet expected criteria.
	T_NotFound               // NotFound errors are returned when a requested resource cannot be located.
	T_Conflict               // Conflict errors occur when a resource already exists.

	CodeInternal = "INTERNAL"
)

// New creates a new ErrorX with the given type, code, and message.
func New(errType Type, code string, msg string) *ErrorX {
	return &ErrorX{
		Code:    code,
		Message: msg,
		Type:    errType,
		origin:  errors.New(msg),
	}
}

// ErrorX represents a structured error that can be used within an application.
type ErrorX struct { // nolint

	// Code is a unique identifier for the error, designed for programmatic consumption.
	Code string `json:"code"`

	// Message provides a human-readable description of the error.
	Message string `json:"message"`

	// Type indicates the category of the error (e.g. Internal, Validation, NotFound).
	Type Type `json:"-"`

	// Details is an optional map of additional details about the error,
	// which is useful for debugging or providing more context in the API responses.
	// For example, it can be used to store the field names that
	// failed validation in the format {"field": "reason"}.
	Details M `json:"details,omitempty"`

	// origin is useful to support error comparison using errors.Is.
	origin error

	// trace stores the trace of the error.
	trace string
}

// Error implements the error interface for ErrorX.
func (e ErrorX) Error() string {
	return e.Message
}

// Trace returns the trace of the error.
func (e *ErrorX) Trace() string {
	return e.trace
}

// WithDetail returns a copy of the ErrorX with an added details.
func (e *ErrorX) WithDetails(details M) *ErrorX {
	newErr := *e
	if newErr.Details == nil {
		newErr.Details = make(map[string]string)
	}
	for k, v := range details {
		newErr.Details[k] = v
	}

	return &newErr
}

// Is implements the errors.Is interface for ErrorX.
// It allows comparison of two ErrorX instances, returning true if they share the same origin.
func (e *ErrorX) Is(target error) bool {
	t, ok := target.(*ErrorX)
	if !ok {
		return false
	}
	return e.origin == t.origin
}

// GetCode returns the error code of an ErrorX.
// If the error is not an ErrorX, it is considered that
// the error is not properly handled in lower layers and returns code Internal.
func GetCode(err error) string {
	if e, ok := err.(*ErrorX); ok {
		return e.Code
	}

	return CodeInternal
}
