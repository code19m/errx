package errx_test

import (
	"strings"
	"testing"

	"github.com/code19m/errx"
)

func TestWithCode(t *testing.T) {
	t.Run("with code", func(t *testing.T) {
		err := errx.New("error", errx.WithCode("1234"))
		e := err.(errx.ErrorX)
		if e.Code() != "1234" {
			t.Errorf("expected code 1234, got %v", e.Code())
		}
	})
}

func TestWithType(t *testing.T) {
	t.Run("with type", func(t *testing.T) {
		err := errx.New("error", errx.WithType(errx.T_Validation))
		e := err.(errx.ErrorX)
		if e.Type() != errx.T_Validation {
			t.Errorf("expected type T_Validation, got %v", e.Type())
		}
	})
}

func TestWithTracePrefix(t *testing.T) {
	t.Run("add prefix to error details and trace", func(t *testing.T) {
		err := errx.New("error", errx.WithDetails(errx.D{"key": "value"}), errx.WithTracePrefix("SERVICE"))
		e := err.(errx.ErrorX)
		if !strings.HasPrefix(e.Trace(), ">>> SERVICE >>>") {
			t.Errorf("expected trace to start with prefix, got: %v", e.Trace())
		}
		if _, ok := e.Details()["SERVICE.key"]; !ok {
			t.Errorf("expected details to include prefixed key")
		}
	})

	t.Run("add prefix to error without details", func(t *testing.T) {
		err := errx.New("error", errx.WithTracePrefix("SERVICE"))
		e := err.(errx.ErrorX)
		if !strings.HasPrefix(e.Trace(), ">>> SERVICE >>>") {
			t.Errorf("expected trace to start with prefix, got: %v", e.Trace())
		}
	})
}

func TestWithDetails(t *testing.T) {
	t.Run("merge new string details with existing ones", func(t *testing.T) {
		err := errx.New("error", errx.WithDetails(errx.D{"key": "value"}))
		err = errx.Wrap(err, errx.WithDetails(errx.D{"key": "new_value"}))
		e := err.(errx.ErrorX)
		if e.Details()["key"] != "new_value | value" {
			t.Errorf("expected merged details, got: %v", e.Details()["key"])
		}
	})

	t.Run("non-string details are replaced", func(t *testing.T) {
		err := errx.New("error", errx.WithDetails(errx.D{"key": 123}))
		err = errx.Wrap(err, errx.WithDetails(errx.D{"key": 456}))
		e := err.(errx.ErrorX)
		if e.Details()["key"] != 456 {
			t.Errorf("expected replaced details, got: %v", e.Details()["key"])
		}
	})

	t.Run("add new details to nil details", func(t *testing.T) {
		// Create an error that doesn't have any details yet
		err := errx.New("error")
		// Add details to it
		err = errx.Wrap(err, errx.WithDetails(errx.D{"new": "detail"}))
		e := err.(errx.ErrorX)
		if e.Details()["new"] != "detail" {
			t.Errorf("expected new detail to be added, got: %v", e.Details())
		}
	})
}

func TestWithFields(t *testing.T) {
	t.Run("overwrite fields", func(t *testing.T) {
		err := errx.New("error", errx.WithFields(errx.M{"field1": "error1"}))
		err = errx.Wrap(err, errx.WithFields(errx.M{"field1": "error2"}))
		e := err.(errx.ErrorX)
		if e.Fields()["field1"] != "error2" {
			t.Errorf("expected overwritten fields, got: %v", e.Fields()["field1"])
		}
	})
}
