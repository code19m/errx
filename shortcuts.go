package errx

// Wr is a shortcut alias for Wrap function.
//
// Wrap wraps an error in an errorX instance with the given options.
//
// This function serves as a convenience wrapper around New,
// enriching the error with additional information and a trace.
//
// It is designed to be used in the middle layers of an application.
var Wr = Wrap

// WrWt is a shortcut alias for WrapWithTypeOnCodes function

// WrapWithTypeOnCodes wraps the error with the given type if the error's code
// matches any of the provided codes. If the error's code doesn't match any of
// the codes, the error type is set to T_Internal.
//
// This function is useful for changing the error type based on its code.
// For example, you can convert an internal error to a validation error
// if its code matches a known validation error code.
//
// If the error is nil, nil is returned.
var WrWt = WrapWithTypeOnCodes
