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

func TestNewf(t *testing.T) {
	t.Run("create new formatted error", func(t *testing.T) {
		err := errx.Newf("error code: %d", 404)
		e, ok := err.(errx.ErrorX)
		if !ok {
			t.Errorf("expected *errx.ErrorX, got %T", err)
		}
		expectedMsg := "error code: 404"
		if e.Error() != expectedMsg {
			t.Errorf("unexpected error message: got %v, want %v", e.Error(), expectedMsg)
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

	t.Run("is returns false for different errors", func(t *testing.T) {
		err1 := fmt.Errorf("error 1")
		wrappedErr := errx.Wrap(err1)
		err2 := fmt.Errorf("error 2")

		if errors.Is(wrappedErr, err2) {
			t.Errorf("expected false for different errors")
		}
	})

	t.Run("is returns false for nil target", func(t *testing.T) {
		err := errx.New("test error")
		if errors.Is(err, nil) {
			t.Errorf("expected false for nil target")
		}
	})

	t.Run("is maintains error identity through wrapping", func(t *testing.T) {
		// Create a base error
		baseErr := fmt.Errorf("base error")

		// Wrap it twice
		wrappedOnce := errx.Wrap(baseErr)
		wrappedTwice := errx.Wrap(wrappedOnce)

		// Check identity is preserved through multiple wrappings
		if !errors.Is(wrappedTwice, baseErr) {
			t.Errorf("expected wrapped error to be identifiable as original through multiple wrappings")
		}
	})

	t.Run("errx and standard errors work together", func(t *testing.T) {
		// Create a base error
		baseErr := fmt.Errorf("standard error")

		// Wrap with errx
		errxWrapped := errx.Wrap(baseErr)

		// Should identify as the original error
		if !errors.Is(errxWrapped, baseErr) {
			t.Errorf("expected errx wrapped error to preserve standard error identity")
		}
	})

	t.Run("errx self-identity check", func(t *testing.T) {
		// Create an ErrorX
		err := errx.New("test error")

		// Check if it identifies as itself
		if !errors.Is(err, err) {
			t.Errorf("error should identify as itself")
		}
	})
}
