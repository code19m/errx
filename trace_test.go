package errx_test

import (
	"strings"
	"testing"

	"github.com/code19m/errx"
)

func TestAddTrace(t *testing.T) {
	t.Run("add trace to error", func(t *testing.T) {
		err := errx.New("error", errx.WithTrace())
		e := err.(errx.ErrorX)
		if e.Trace() == "" {
			t.Errorf("expected trace to be populated, got empty")
		}
	})

	t.Run("verify trace format", func(t *testing.T) {
		err := errx.New("error", errx.WithTrace())
		e := err.(errx.ErrorX)
		if !contains(e.Trace(), "trace_test.go") {
			t.Errorf("expected trace to include filename, got: %v", e.Trace())
		}
	})
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
