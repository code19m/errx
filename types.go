package errx

import "fmt"

const (
	// Internal errors indicate unexpected issues within the application.
	T_Internal Type = iota

	// Validation errors occur when user input does not meet expected criteria.
	T_Validation

	// NotFound errors are returned when a requested resource cannot be located.
	T_NotFound

	// Conflict errors occur when a resource already exists.
	T_Conflict

	// Authentication errors occur when a user is not authorized to access a resource.
	T_Authentication

	// Forbidden errors occur when a user is not allowed to access a resource.
	T_Forbidden
)

const (
	// DefaultCode is the default error code used when no code is provided.
	DefaultCode = ""

	// DefaultType is the default error type used when no type is provided.
	DefaultType = T_Internal
)

// Type defines the different categories of errors that can be represented by an ErrorX.
type Type uint8

// M is a shorthand for a map of string key-value pairs.
type M map[string]string

// String returns a string representation of the error type.
func (t Type) String() string {
	switch t {
	case T_Internal:
		return "T_Internal"
	case T_Validation:
		return "T_Validation"
	case T_NotFound:
		return "T_NotFound"
	case T_Conflict:
		return "T_Conflict"
	case T_Authentication:
		return "T_Authentication"
	case T_Forbidden:
		return "T_Forbidden"
	default:
		return fmt.Sprintf("Unknown Type (%d)", t)
	}
}
