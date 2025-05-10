package errx_test

import (
	"errors"
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
		if st.Message() != "internal server error" {
			t.Errorf("unexpected GRPC error message: %v", st.Message())
		}
	})

	t.Run("convert nil error to nil", func(t *testing.T) {
		if err := errx.ToGRPCError(nil); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("convert regular error to GRPC error", func(t *testing.T) {
		err := errors.New("regular error")
		grpcErr := errx.ToGRPCError(err)
		st, ok := status.FromError(grpcErr)
		if !ok {
			t.Errorf("expected GRPC status error")
		}
		if st.Message() != "regular error" {
			t.Errorf("unexpected GRPC error message: %v", st.Message())
		}
	})

	t.Run("apply options when converting to GRPC error", func(t *testing.T) {
		err := errors.New("error with options")
		grpcErr := errx.ToGRPCError(err, errx.WithCode("CUSTOM_CODE"), errx.WithType(errx.T_Validation))

		// Extract the status and details
		st, ok := status.FromError(grpcErr)
		if !ok {
			t.Errorf("expected GRPC status error")
		}

		// Check that the code in the proto details is correct
		if len(st.Details()) == 0 {
			t.Errorf("expected GRPC error to have details")
		}

		// Verify the status code maps to validation
		if st.Code() != codes.InvalidArgument {
			t.Errorf("expected code %v, got %v", codes.InvalidArgument, st.Code())
		}
	})

	t.Run("convert ErrorX with different types to proper GRPC codes", func(t *testing.T) {
		testCases := []struct {
			errType  errx.Type
			expected codes.Code
		}{
			{errx.T_Internal, codes.Internal},
			{errx.T_Validation, codes.InvalidArgument},
			{errx.T_NotFound, codes.NotFound},
			{errx.T_Conflict, codes.AlreadyExists},
			{errx.T_Authentication, codes.Unauthenticated},
			{errx.T_Forbidden, codes.PermissionDenied},
			{errx.Type(99), codes.Unknown}, // Unknown type should map to unknown code
		}

		for _, tc := range testCases {
			err := errx.New("test error", errx.WithType(tc.errType))
			grpcErr := errx.ToGRPCError(err)
			st, _ := status.FromError(grpcErr)

			if st.Code() != tc.expected {
				t.Errorf("for error type %v, expected gRPC code %v, got %v",
					tc.errType, tc.expected, st.Code())
			}
		}
	})
}

func TestFromGRPCError(t *testing.T) {
	t.Run("convert gRPC error to ErrorX", func(t *testing.T) {
		grpcErr := errx.ToGRPCError(errx.New("not found", errx.WithCode("NOT_FOUND"), errx.WithType(errx.T_NotFound)))

		ok, err := errx.FromGRPCError(grpcErr)
		e := err.(errx.ErrorX)

		if !ok {
			t.Errorf("expected successful conversion")
		}
		if e.Type() != errx.T_NotFound {
			t.Errorf("expected type T_NotFound, got %v", e.Type())
		}
		if e.Code() != "NOT_FOUND" {
			t.Errorf("expected code %v, got %v", codes.NotFound.String(), e.Code())
		}
	})

	t.Run("convert gRPC error with no ErrorX detail", func(t *testing.T) {
		grpcErr := status.Error(codes.AlreadyExists, "resource already exists")

		ok, err := errx.FromGRPCError(grpcErr)
		e := err.(errx.ErrorX)

		if ok {
			t.Errorf("expected unsuccessful conversion")
		}
		if e.Type() != errx.T_Conflict {
			t.Errorf("expected type T_AlreadyExists, got %v", e.Type())
		}
		if e.Code() != errx.DefaultCode {
			t.Errorf("expected default code, got %v", e.Code())
		}
	})

	t.Run("test all gRPC code mappings to error types", func(t *testing.T) {
		testCases := []struct {
			code     codes.Code
			expected errx.Type
		}{
			{codes.Internal, errx.T_Internal},
			{codes.InvalidArgument, errx.T_Validation},
			{codes.NotFound, errx.T_NotFound},
			{codes.AlreadyExists, errx.T_Conflict},
			{codes.Unauthenticated, errx.T_Authentication},
			{codes.PermissionDenied, errx.T_Forbidden},
			{codes.Unknown, errx.T_Internal}, // Unknown should map to Internal
		}

		for _, tc := range testCases {
			grpcErr := status.Error(tc.code, "test error")
			ok, err := errx.FromGRPCError(grpcErr)
			e := err.(errx.ErrorX)

			if ok {
				t.Errorf("expected unsuccessful conversion for code %v", tc.code)
			}
			if e.Type() != tc.expected {
				t.Errorf("for gRPC code %v, expected error type %v, got %v",
					tc.code, tc.expected, e.Type())
			}
		}
	})

	t.Run("handle non-gRPC error", func(t *testing.T) {
		genericErr := errx.New("generic error")
		ok, err := errx.FromGRPCError(genericErr)
		e := err.(errx.ErrorX)

		if ok {
			t.Errorf("expected unsuccessful conversion")
		}
		if e.Type() != errx.T_Internal {
			t.Errorf("expected type T_Internal, got %v", e.Type())
		}
		if e.Code() != errx.DefaultCode {
			t.Errorf("expected default code, got %v", e.Code())
		}
	})

	t.Run("handle nil gRPC error", func(t *testing.T) {
		ok, err := errx.FromGRPCError(nil)
		if ok {
			t.Errorf("expected unsuccessful conversion")
		}
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
