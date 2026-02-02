package errx

import "slices"

// GetCode returns the error code if the error implements the ErrorX interface.
// Otherwise, it returns the default code.
func GetCode(err error) string {
	if e, ok := err.(ErrorX); ok {
		return e.Code()
	}
	return DefaultCode
}

// GetType returns the error type if the error implements the ErrorX interface.
// Otherwise, it returns the default type.
func GetType(err error) Type {
	if e, ok := err.(ErrorX); ok {
		return e.Type()
	}
	return DefaultType
}

// IsCodeIn checks if the error's code is in the given list of codes.
func IsCodeIn(err error, codes ...string) bool {
	code := GetCode(err)
	return slices.Contains(codes, code)
}

// AsErrorX returns the error as an ErrorX instance.
//
// If the error does not implement the ErrorX interface,
// it converts it to an ErrorX with default values.
//
// This function is useful when you want to work with ErrorX instances.
//
// ***NOTE***: Make sure that error is not nil before calling this function.
func AsErrorX(err error) ErrorX {
	if e, ok := err.(ErrorX); ok {
		return e
	}

	return wrapFromError(err)
}

// WrapWithTypeOnCodes wraps the error with the given type if the error's code
// matches any of the provided codes. If the error's code doesn't match any of
// the codes, the error type remains unchanged.
//
// This function is useful for changing the error type based on its code.
// For example, you can convert an internal error to a validation error
// if its code matches a known validation error code.
//
// If the error is nil, nil is returned.
func WrapWithTypeOnCodes(err error, type_ Type, codes ...string) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*errorX)
	if !ok {
		e = wrapFromError(err)
	}

	if slices.Contains(codes, e.Code()) {
		e.type_ = type_
	}

	e = e.clone()
	e.addTrace(2)
	return e
}
