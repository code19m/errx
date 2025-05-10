package errx

import (
	"errors"
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
		e = newDefault(err.Error())
	}

	// Clone the error to avoid modifying the original
	e = e.clone()

	// Apply options
	e.addTrace(2)
	applyOpts(e, opts)

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
// It is intended for use on the gRPC client side after making gRPC calls, to convert gRPC status errors into ErrorX instances.
// The function returns a boolean indicating whether the error was successfully converted from a proto message (`true`) or not (`false`).
//
// If the error is nil, no action is taken, and the function returns `false, nil`.
// If the provided error does not contain an ErrorX detail, a default ErrorX instance is created.
// Optional modifications can be applied via OptionFunc.
//
// ***NOTE***: Don't confuse this function with ToGRPCError, which is intended for use on the gRPC server side.
func FromGRPCError(err error, opts ...OptionFunc) (bool, error) {
	if err == nil {
		return false, nil
	}

	st, ok := status.FromError(err)
	if !ok {
		e := newDefault(err.Error())
		e.addTrace(2)
		applyOpts(e, opts)
		return false, e
	}

	for _, detail := range st.Details() {
		if pb, ok := detail.(*errorx_proto.ErrorX); ok {
			e := fromProto(pb)
			e.addTrace(2)
			applyOpts(e, opts)
			return true, e
		}
	}

	e := newFromStatus(st)
	e.addTrace(2)
	applyOpts(e, opts)
	return false, e
}

// fromProto converts a proto error to an ErrorX.
func fromProto(pbErr *errorx_proto.ErrorX) *errorX {
	return &errorX{
		code:    pbErr.GetCode(),
		msg:     pbErr.GetMessage(),
		type_:   Type(pbErr.GetType()),
		fields:  M(pbErr.GetFields()),
		details: make(D),
		trace:   pbErr.GetTrace(),
		origin:  errors.New(pbErr.GetMessage()),
	}
}

// toProto converts an ErrorX to a proto error.
func toProto(e *errorX) *errorx_proto.ErrorX {
	return &errorx_proto.ErrorX{
		Code:    e.Code(),
		Message: e.Error(),
		Type:    int32(e.Type()),
		Fields:  e.Fields(),
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

// newFromStatus creates a new ErrorX from a gRPC status.
// This function is used when the gRPC status does not contain an ErrorX in its details.
func newFromStatus(st *status.Status) *errorX {
	rpcCodeMap := map[codes.Code]Type{
		codes.Internal:         T_Internal,
		codes.InvalidArgument:  T_Validation,
		codes.NotFound:         T_NotFound,
		codes.AlreadyExists:    T_Conflict,
		codes.Unauthenticated:  T_Authentication,
		codes.PermissionDenied: T_Forbidden,
	}

	if t, ok := rpcCodeMap[st.Code()]; ok {
		return &errorX{
			code:    DefaultCode,
			msg:     st.Message(),
			type_:   t,
			fields:  make(M),
			details: make(D),
		}
	}

	return newDefault(st.String())
}
