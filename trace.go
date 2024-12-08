package errx

import (
	"fmt"
	"runtime"
	"strings"
)

// addTrace captures and records the caller's context at a specific point in the error's propagation.
// It collects detailed information about the function that triggered the error trace,
// helping developers understand the exact path of error generation.
//
// Parameters:
//   - skipNumber: Controls the depth of stack trace collection
//   - 0: current function
//   - 1: immediate caller
//   - 2: caller of the function that invoked addTrace
//   - 3: caller of the function that invoked the function that invoked addTrace
//
// The method builds a compact, human-readable trace that includes:
//   - Filename
//   - Line number
//   - Function name (without full package path)
//
// The trace is constructed in chronological order,
// allowing developers to see the most recent call first.
func (e *errorX) addTrace(skipNumber int) {
	// Retrieve caller information using runtime reflection
	// Panics if unable to obtain caller details to prevent silent failures
	pc, filepath, line, ok := runtime.Caller(skipNumber)
	if !ok {
		panic("could not get runtime.Caller")
	}

	// Retrieve function details for the program counter
	// Ensures we have valid function metadata for tracing
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		panic("could not get runtime.FuncForPC")
	}

	// Extract filename without full path to keep trace compact
	_, filename := pathSplit(filepath)

	// Get function name, removing package prefix for readability
	funcName := fn.Name()
	shortFuncName := funcName[strings.LastIndex(funcName, "/")+1:]

	// Format caller information with filename, line, and function name
	callerInfo := fmt.Sprintf("[%s:%d] %s", filename, line, shortFuncName)

	// Prepend new trace information, creating a chain of function calls
	// Uses right-pointing arrow (➡️) to visually represent call progression
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
