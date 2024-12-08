package errx

import (
	"fmt"

	"github.com/code19m/errx/internal/errorx_proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError converts a custom error (ErrorX) into a gRPC-compatible error.
//
// It is intended for use in gRPC server handlers/interceptors to convert ErrorX instances to gRPC status errors.
// If the error is nil, no action is taken, so it is safe to call this function with a nil error.
//
// If the provided error does not implement the ErrorX interface, it is wrapped
// into a default ErrorX instance.
//
// Optional modifications can be applied via OptionFunc.
//
// ***NOTE***: Don't confuse this function with FromGRPCError, which is intended for use in gRPC client side.
func ToGRPCError(err error, opts ...OptionFunc) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*errorX)
	if !ok {
		e = New(err.Error()).(*errorX)
	}

	// Apply options
	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}

	st, derr := status.New(mapErrorToGRPCCode(e), e.Error()).WithDetails(toProto(e))
	if derr != nil {
		return New(
			fmt.Sprintf(
				"Failed to create grpc status object with details: %s. Original error was: %s",
				derr.Error(),
				e.Error(),
			),
		)
	}

	return st.Err()
}

// FromGRPCError converts a gRPC error into a custom error (ErrorX).
//
// It is intended for use in gRPC client side after making gRPC calls to convert gRPC status errors to ErrorX instances.
// If the error is nil, no action is taken, so it is safe to call this function with a nil error.
//
// If the provided error does not contain an ErrorX detail, a default ErrorX instance is created.
// Optional modifications can be applied via OptionFunc.
//
// ***NOTE***: Don't confuse this function with ToGRPCError, which is intended for use in gRPC server side.
func FromGRPCError(err error, opts ...OptionFunc) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return New(err.Error(), opts...)
	}

	for _, detail := range st.Details() {
		if pb, ok := detail.(*errorx_proto.ErrorX); ok {
			e := fromProto(pb)

			// Apply options
			for _, opt := range opts {
				if opt != nil {
					opt(e)
				}
			}

			return e
		}
	}

	return New(err.Error(), opts...)
}

// fromProto converts a proto error to an ErrorX.
func fromProto(pbErr *errorx_proto.ErrorX) *errorX {
	return &errorX{
		code:    pbErr.GetCode(),
		msg:     pbErr.GetMessage(),
		type_:   Type(pbErr.GetType()),
		fields:  M(pbErr.GetFields()),
		details: M(pbErr.GetDetails()),
		trace:   pbErr.GetTrace(),
	}
}

// toProto converts an ErrorX to a proto error.
func toProto(e *errorX) *errorx_proto.ErrorX {
	return &errorx_proto.ErrorX{
		Code:    e.Code(),
		Message: e.Error(),
		Type:    int32(e.Type()),
		Fields:  e.Fields(),
		Details: e.Details(),
		Trace:   e.Trace(),
	}
}

// mapErrorToGRPCCode returns the gRPC code for an ErrorX based on its type.
func mapErrorToGRPCCode(err *errorX) codes.Code {
	switch err.Type() {
	case T_Internal:
		return codes.Internal
	case T_Validation:
		return codes.InvalidArgument
	case T_NotFound:
		return codes.NotFound
	case T_Conflict:
		return codes.AlreadyExists
	case T_Authentication:
		return codes.Unauthenticated
	case T_Forbidden:
		return codes.PermissionDenied
	}
	return codes.Unknown
}
