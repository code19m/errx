package errx_test

import (
	"testing"

	"github.com/code19m/errx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPCError(t *testing.T) {
	t.Run("convert ErrorX to GRPC error", func(t *testing.T) {
		err := errx.New("internal server error", errx.WithCode("INTERNAL"))
		grpcErr := errx.ToGRPCError(err)
		st, ok := status.FromError(grpcErr)
		if !ok {
			t.Errorf("expected GRPC status error")
		}
		if st.Message() != "[T_Internal: INTERNAL] - internal server error" {
			t.Errorf("unexpected GRPC error message: %v", st.Message())
		}
	})
}

func TestFromGRPCError(t *testing.T) {
	t.Run("convert gRPC error to ErrorX", func(t *testing.T) {
		grpcErr := errx.ToGRPCError(errx.New("not found", errx.WithCode("NOT_FOUND"), errx.WithType(errx.T_NotFound)))

		err := errx.FromGRPCError(grpcErr)
		e := err.(errx.ErrorX)

		if e.Type() != errx.T_NotFound {
			t.Errorf("expected type T_NotFound, got %v", e.Type())
		}
		if e.Code() != "NOT_FOUND" {
			t.Errorf("expected code %v, got %v", codes.NotFound.String(), e.Code())
		}
	})

	t.Run("convert gRPC error with no ErrorX detail", func(t *testing.T) {
		grpcErr := status.Error(codes.AlreadyExists, "resource already exists")

		err := errx.FromGRPCError(grpcErr)
		e := err.(errx.ErrorX)

		if e.Type() != errx.T_Conflict {
			t.Errorf("expected type T_AlreadyExists, got %v", e.Type())
		}
		if e.Code() != errx.DefaultCode {
			t.Errorf("expected default code, got %v", e.Code())
		}
	})

	t.Run("handle non-gRPC error", func(t *testing.T) {
		genericErr := errx.New("generic error")
		err := errx.FromGRPCError(genericErr)
		e := err.(errx.ErrorX)

		if e.Type() != errx.T_Internal {
			t.Errorf("expected type T_Internal, got %v", e.Type())
		}
		if e.Code() != errx.DefaultCode {
			t.Errorf("expected default code, got %v", e.Code())
		}
	})

	t.Run("handle nil gRPC error", func(t *testing.T) {
		err := errx.FromGRPCError(nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
