package errx_test

import (
	"fmt"
	"testing"

	"github.com/code19m/errx"
)

func TestGetCode(t *testing.T) {
	t.Run("get code from ErrorX", func(t *testing.T) {
		err := errx.New("error", errx.WithCode("CUSTOM_CODE"))
		if code := errx.GetCode(err); code != "CUSTOM_CODE" {
			t.Errorf("expected code CUSTOM_CODE, got %v", code)
		}
	})

	t.Run("get default code for non-ErrorX error", func(t *testing.T) {
		err := fmt.Errorf("generic error")
		if code := errx.GetCode(err); code != errx.DefaultCode {
			t.Errorf("expected default code, got %v", code)
		}
	})
}

func TestGetType(t *testing.T) {
	t.Run("get type from ErrorX", func(t *testing.T) {
		err := errx.New("error", errx.WithType(errx.T_Validation))
		if typ := errx.GetType(err); typ != errx.T_Validation {
			t.Errorf("expected type T_Validation, got %v", typ)
		}
	})

	t.Run("get default type for non-ErrorX error", func(t *testing.T) {
		err := fmt.Errorf("generic error")
		if typ := errx.GetType(err); typ != errx.DefaultType {
			t.Errorf("expected default type, got %v", typ)
		}
	})
}

func TestAsErrorX(t *testing.T) {
	t.Run("convert to ErrorX from generic error", func(t *testing.T) {
		err := fmt.Errorf("generic error")
		e := errx.AsErrorX(err)
		if e.Error() != "generic error" {
			t.Errorf("unexpected ErrorX message: %v", e.Error())
		}
	})
}
