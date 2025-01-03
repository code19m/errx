package errx

import (
	"errors"
)

// ErrorX represents a main interface of this package.
// It extends the built-in error interface with additional methods
// to provide structured error information and facilitate debugging.
type ErrorX interface {

	// Error returns a human-readable description of the error.
	// It implements the standard error interface.
	Error() string

	// Code returns a machine-readable error code.
	// This is intended for use in application logic.
	Code() string

	// Type returns the type of the error.
	// Useful for categorizing errors during error handling.
	Type() Type

	// Trace returns the error's trace information.
	// This can help identify the error's origin in the system.
	Trace() string

	// Fields provides information about input validation errors.
	// Example: {"field_name": "error_message/validation_rule"}
	// Not to be confused with Details, which is used for debugging.
	Fields() M

	// Details provides additional debugging information about the error.
	// This is intended for logging and troubleshooting purposes.
	Details() M

	// Is methos implements the standard errors.Is function.
	// It reports whether any error in the error's tree matches the target.
	Is(target error) bool
}

// New creates a new ErrorX with the given message and options.
func New(msg string, opts ...OptionFunc) error {
	e := newDefault(msg)

	// Apply options
	e.addTrace(2)
	applyOpts(e, opts)

	return e
}

// Wrap wraps an error in an errorX instance with the given options.
//
// This function serves as a convenience wrapper around New,
// enriching the error with additional information and a trace.
//
// It is designed to be used in the middle layers of an application.
func Wrap(err error, opts ...OptionFunc) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*errorX)
	if !ok {
		e = wrapFromError(err)
	}

	// Clone the error to avoid modifying the original
	e = e.clone()

	// Apply options
	e.addTrace(2)
	applyOpts(e, opts)

	return e
}

// errorX is a concrete implementation of the ErrorX interface.
type errorX struct {
	code    string
	msg     string
	type_   Type
	fields  M
	details M
	trace   string
	origin  error
}

func (e errorX) Error() string {
	return e.msg
}

func (e errorX) Code() string {
	return e.code
}

func (e errorX) Type() Type {
	return e.type_
}

func (e errorX) Trace() string {
	return e.trace
}

func (e errorX) Fields() M {
	return e.fields
}

func (e errorX) Details() M {
	return e.details
}

func (e errorX) Is(target error) bool {
	if target == nil {
		return false
	}
	return e.origin == target
}

func (e *errorX) clone() *errorX {
	return &errorX{
		code:    e.code,
		msg:     e.msg,
		type_:   e.type_,
		fields:  e.fields,
		details: e.details,
		trace:   e.trace,
		origin:  e.origin,
	}
}

func newDefault(msg string) *errorX {
	return &errorX{
		code:    DefaultCode,
		msg:     msg,
		type_:   DefaultType,
		fields:  make(M),
		details: make(M),
		origin:  errors.New(msg),
	}
}

func wrapFromError(err error) *errorX {
	return &errorX{
		code:    DefaultCode,
		msg:     err.Error(),
		type_:   DefaultType,
		fields:  make(M),
		details: make(M),
		origin:  err,
	}
}

func applyOpts(e *errorX, opts []OptionFunc) {
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}
}
