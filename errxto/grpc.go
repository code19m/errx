package errxto

import (
	"fmt"

	"github.com/code19m/errx"
	"github.com/code19m/errx/internal/errpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errto.GRPC is intended for use in gRPC server handlers to convert ErrorX instances to gRPC status errors.
// If the error is nil, no response is written.
func GRPC(err error) error {
	if err == nil {
		return nil
	}

	_, ok := err.(*errx.ErrorX)
	if !ok {
		err = errx.Wrap(err)
	}

	return toStatus(err).Err()
}

func toStatus(err error) *status.Status {
	if e, ok := err.(*errx.ErrorX); ok {
		st, dtErr := status.New(gRPCStatusCode(e), e.Message).WithDetails(
			&errpb.ErrorX{
				Message: e.Message,
				Code:    e.Code,
				Type:    int32(e.Type),
				Details: e.Details,
				Trace:   e.Trace(),
			},
		)
		if dtErr == nil {
			return st
		}
		err = fmt.Errorf("st.WithDetails error: %w original error: %w", dtErr, err)
	}
	return status.New(codes.Internal, err.Error())
}

func gRPCStatusCode(err error) codes.Code {
	if e, ok := err.(*errx.ErrorX); ok {
		switch e.Type {
		case errx.T_Validation:
			return codes.InvalidArgument
		case errx.T_NotFound:
			return codes.NotFound
		case errx.T_Conflict:
			return codes.AlreadyExists
		case errx.T_Internal:
			return codes.Internal
		}
	}
	return codes.Internal
}
