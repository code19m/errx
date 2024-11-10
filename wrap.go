package errx

import (
	"fmt"
	"runtime"
	"strings"
)

// Wrap wraps an error with an ErrorX.
// The function appends the caller information as
// a stacktrace of the function that invoked the function
// If the error is nil, Wrap returns nil.
// If the error is not an instance of ErrorX, then created New ErrorX with
// type Internal and the err.Error() as a message.
// You don't need to call Wrap after functions: New, AsInternal, FromPgxQuery, FromPgxExec, FromGRPC
// because they already add a trace to the error.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*ErrorX)
	if !ok {
		e = &ErrorX{
			Code:    CodeInternal,
			Message: err.Error(),
			Type:    T_Internal,
			origin:  err,
		}
	}

	e.addTrace()
	return e
}

// WrapWithDetails wraps an error with an ErrorX and adds details to it.
func WrapWithDetails(err error, details M) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*ErrorX)
	if !ok {
		e = &ErrorX{
			Code:    CodeInternal,
			Message: err.Error(),
			Type:    T_Internal,
			origin:  err,
		}
	}

	e.addTrace()
	return e.WithDetails(details)
}

// addTrace appends the caller information of the function that invoked the function
// that called addTrace to the error's trace field. This helps in tracking the chain
// of function calls leading to the error, providing a detailed trace for debugging.
func (e *ErrorX) addTrace() {
	// Skip 0 to get the current function, 1 to get the caller of the current function, etc.
	// Here, we skip 2 to get the caller of the function that invoked addTrace.
	pc, filepath, line, ok := runtime.Caller(2)
	if !ok {
		panic("could not get runtime.Caller")
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		panic("could not get runtime.FuncForPC")
	}

	// Shorten the file path to make it more compact
	_, filename := pathSplit(filepath)

	// Get the function name (this can be simplified)
	funcName := fn.Name()

	// Get the short function name (without the package path)
	shortFuncName := funcName[strings.LastIndex(funcName, ".")+1:]

	callerInfo := fmt.Sprintf("[%s:%d] %s", filename, line, shortFuncName)

	if e.trace == "" {
		e.trace = callerInfo
	} else {
		e.trace = fmt.Sprintf("%s ➡️ %s", callerInfo, e.trace)
	}
}

// pathSplit splits a path into the directory and the file name
func pathSplit(path string) (string, string) {
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}
	return "", path
}
