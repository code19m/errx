package errx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/code19m/errx"
)

func TestNew(t *testing.T) {
	t.Run("create new error with default values", func(t *testing.T) {
		err := errx.New("test error")
		e, ok := err.(errx.ErrorX)
		if !ok {
			t.Errorf("expected *errx.ErrorX, got %T", err)
		}
		if e.Error() != "test error" {
			t.Errorf("unexpected error message: %v", e.Error())
		}
	})

	t.Run("create new error with custom code and type", func(t *testing.T) {
		err := errx.New("validation failed", errx.WithCode("VALIDATION_ERROR"), errx.WithType(errx.T_Validation))
		e := err.(errx.ErrorX)
		if e.Code() != "VALIDATION_ERROR" || e.Type() != errx.T_Validation {
			t.Errorf("unexpected code or type: %v, %v", e.Code(), e.Type())
		}
	})
}

func TestWrap(t *testing.T) {
	t.Run("wrap nil error", func(t *testing.T) {
		err := errx.Wrap(nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("wrap non-errorX error", func(t *testing.T) {
		originalErr := fmt.Errorf("an external error")
		err := errx.Wrap(originalErr)
		e := err.(errx.ErrorX)
		if e.Error() != "an external error" {
			t.Errorf("unexpected wrapped error message: %v", e.Error())
		}
	})
}

func TestIs(t *testing.T) {
	t.Run("check if error is wrapped", func(t *testing.T) {
		originalErr := fmt.Errorf("an external error")
		wrappedErr := errx.Wrap(originalErr)
		if !errors.Is(wrappedErr, originalErr) {
			t.Errorf("expected true, got false")
		}
	})
}
