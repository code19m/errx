package errx

import "fmt"

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
}

// New creates a new ErrorX with the given message and options.
func New(msg string, opts ...OptionFunc) error {
	e := &errorX{
		code:  DefaultCode,
		msg:   msg,
		type_: DefaultType,
	}

	// Apply options
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}

	return e
}

// Wrap wraps an error in an errorX instance with the given options.
//
// This function serves as a convenience wrapper around New,
// enriching the error with additional information and a trace.
//
// It is designed to be used in the middle layers of an application.
//
// ***NOTE***: Do not use the WithTrace option with this function,
// as a trace is added by default when using Wrap.
func Wrap(err error, opts ...OptionFunc) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*errorX)
	if !ok {
		e = &errorX{
			code:  DefaultCode,
			msg:   err.Error(),
			type_: DefaultType,
		}
	}
	e.addTrace(2)

	// Apply options
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}

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
}

func (e errorX) Error() string {
	return fmt.Sprintf("[%s: %s] - %s", e.type_, e.code, e.msg)
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
