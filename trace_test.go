package errx_test

import (
	"strings"
	"testing"

	"github.com/code19m/errx"
)

func TestAddTrace(t *testing.T) {
	t.Run("add trace to error", func(t *testing.T) {
		err := errx.New("error")
		e := err.(errx.ErrorX)
		if e.Trace() == "" {
			t.Errorf("expected trace to be populated, got empty")
		}
	})

	t.Run("verify trace format", func(t *testing.T) {
		err := errx.New("error")
		e := err.(errx.ErrorX)
		if !contains(e.Trace(), "trace_test.go") {
			t.Errorf("expected trace to include filename, got: %v", e.Trace())
		}
	})

	t.Run("test trace propagation", func(t *testing.T) {
		// Create a function chain to test trace propagation
		err := generateErrorThroughChain()
		e := err.(errx.ErrorX)

		// The trace should contain all three functions in the chain
		trace := e.Trace()
		if !contains(trace, "generateErrorThroughChain") {
			t.Errorf("trace missing caller function: %s", trace)
		}
		if !contains(trace, "middleFunction") {
			t.Errorf("trace missing middle function: %s", trace)
		}
		if !contains(trace, "innerErrorGenerator") {
			t.Errorf("trace missing inner function: %s", trace)
		}
	})

	t.Run("test different file paths in trace", func(t *testing.T) {
		// Test different path formats through public API
		errorFromFunc := func() error {
			return errx.New("error from func")
		}
		err := errorFromFunc()
		e := err.(errx.ErrorX)

		// Checking that filename extraction works with different path formats
		if !contains(e.Trace(), "trace_test.go") {
			t.Errorf("expected trace to include filename, got: %v", e.Trace())
		}
	})
}

// Helper functions to create a trace chain
func generateErrorThroughChain() error {
	return errx.Wrap(middleFunction())
}

func middleFunction() error {
	return errx.Wrap(innerErrorGenerator())
}

func innerErrorGenerator() error {
	return errx.New("inner error")
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
