package errx

import "fmt"

// OptionFunc is a function that modifies an errorX.
type OptionFunc func(*errorX)

// WithCode sets the error code.
// If this option is not used, the default code is "INTERNAL".
func WithCode(code string) OptionFunc {
	return func(e *errorX) {
		e.code = code
	}
}

// WithType sets the error type.
// If this option is not used, the default type is T_Internal.
func WithType(t Type) OptionFunc {
	return func(e *errorX) {
		e.type_ = t
	}
}

// WithPrefix adds a prefix to the trace and all keys in the error's details,
// specifically designed for error propagation between microservices,
// particularly in gRPC communication.
//
// The trace is changed in the format ">>> prefix >>> %s".
// The details keys are changed in the format "prefix.%s".
func WithPrefix(prefix string) OptionFunc {
	return func(e *errorX) {
		e.trace = fmt.Sprintf(">>> %s >>> %s", prefix, e.trace)
		if e.details != nil {
			details := make(D)
			for k, v := range e.details {
				details[fmt.Sprintf("%s.%s", prefix, k)] = v
			}
			e.details = details
		}
	}
}

// WithDetails adds additional contextual information (metadata) to the error.
// If a key already exists, the new value is appended to the existing value,
// with the new value appearing first, separated by a "|" character if both values are strings.
func WithDetails(details D) OptionFunc {
	return func(e *errorX) {
		if e.details == nil {
			e.details = make(D)
		}

		for k, v := range details {
			if existingVal, ok := e.details[k]; ok {
				// If both values are strings, concatenate with a separator
				if strVal, isStr := v.(string); isStr {
					if strExistingVal, isExistingStr := existingVal.(string); isExistingStr {
						e.details[k] = strVal + " | " + strExistingVal
						continue
					}
				}
				// Otherwise, just replace with the new value
				e.details[k] = v
			} else {
				e.details[k] = v
			}
		}
	}
}

// WithFields sets specific validation related fields.
// Unlike WithDetails, this method does not append but completely overwrites the existing fields.
//
// Example of fields:
// {"username": "too short", "email": "invalid format"}
func WithFields(fields M) OptionFunc {
	return func(e *errorX) {
		e.fields = fields
	}
}
