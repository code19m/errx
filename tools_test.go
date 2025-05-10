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

	t.Run("convert from ErrorX stays the same", func(t *testing.T) {
		originalErr := errx.New("already an ErrorX")
		e := errx.AsErrorX(originalErr)
		if e.Error() != "already an ErrorX" {
			t.Errorf("unexpected ErrorX message: %v", e.Error())
		}
	})
}

func TestIsCodeIn(t *testing.T) {
	t.Run("code is in the list", func(t *testing.T) {
		err := errx.New("error", errx.WithCode("CODE_1"))
		if !errx.IsCodeIn(err, "CODE_0", "CODE_1", "CODE_2") {
			t.Errorf("expected true for code CODE_1 in the list")
		}
	})

	t.Run("code is not in the list", func(t *testing.T) {
		err := errx.New("error", errx.WithCode("CODE_3"))
		if errx.IsCodeIn(err, "CODE_0", "CODE_1", "CODE_2") {
			t.Errorf("expected false for code CODE_3 not in the list")
		}
	})

	t.Run("default code for non-ErrorX error", func(t *testing.T) {
		err := fmt.Errorf("generic error")
		if !errx.IsCodeIn(err, errx.DefaultCode) {
			t.Errorf("expected true for default code in the list")
		}
	})
}

func TestWrapWithTypeOnCodes(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		var err error
		result := errx.WrapWithTypeOnCodes(err, errx.T_Validation, "CODE_1")
		if result != nil {
			t.Errorf("expected nil for nil error, got %v", result)
		}
	})

	t.Run("wrap with type when code matches", func(t *testing.T) {
		err := errx.New("error", errx.WithCode("CODE_1"), errx.WithType(errx.T_Internal))
		result := errx.WrapWithTypeOnCodes(err, errx.T_Validation, "CODE_1", "CODE_2")
		
		if errx.GetType(result) != errx.T_Validation {
			t.Errorf("expected type T_Validation, got %v", errx.GetType(result))
		}
	})

	t.Run("keep original type when code doesn't match", func(t *testing.T) {
		err := errx.New("error", errx.WithCode("CODE_3"), errx.WithType(errx.T_Internal))
		result := errx.WrapWithTypeOnCodes(err, errx.T_Validation, "CODE_1", "CODE_2")
		
		if errx.GetType(result) != errx.T_Internal {
			t.Errorf("expected type T_Internal, got %v", errx.GetType(result))
		}
	})

	t.Run("wrap non-ErrorX error", func(t *testing.T) {
		err := fmt.Errorf("generic error")
		result := errx.WrapWithTypeOnCodes(err, errx.T_Validation, errx.DefaultCode)
		
		if errx.GetType(result) != errx.T_Validation {
			t.Errorf("expected type T_Validation, got %v", errx.GetType(result))
		}
	})
}
