package errx

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
	return Wrap(err).(ErrorX)
}
