package errx_test

import (
	"testing"

	"github.com/code19m/errx"
)

func TestTypeString(t *testing.T) {
	tests := []struct {
		typ      errx.Type
		expected string
	}{
		{errx.T_Internal, "T_Internal"},
		{errx.T_Validation, "T_Validation"},
		{errx.T_NotFound, "T_NotFound"},
		{errx.T_Conflict, "T_Conflict"},
		{errx.T_Authentication, "T_Authentication"},
		{errx.T_Forbidden, "T_Forbidden"},
		{errx.Type(99), "Unknown Type (99)"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if str := test.typ.String(); str != test.expected {
				t.Errorf("expected %v, got %v", test.expected, str)
			}
		})
	}
}
