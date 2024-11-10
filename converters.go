package errx

import (
	"errors"
	"fmt"

	"github.com/code19m/errx/internal/errpb"
	"google.golang.org/grpc/status"
)

func FromGRPC(err error, serviceName string) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		e := &ErrorX{
			Code:    CodeInternal,
			Message: err.Error(),
			Type:    T_Internal,
			origin:  err,
		}
		e.addTrace()
		return e
	}

	for _, detail := range st.Details() {
		if pb, ok := detail.(*errpb.ErrorX); ok {
			e := fromProto(pb, serviceName)
			e.addTrace()
			return e
		}
	}
	e := &ErrorX{
		Code:    CodeInternal,
		Message: err.Error(),
		Type:    T_Internal,
		origin:  err,
	}
	e = e.WithDetails(M{
		fmt.Sprintf("%s.%s", serviceName, "grpc_code"): st.Code().String(),
		fmt.Sprintf("%s.%s", serviceName, "grpc_msg"):  st.Message(),
	})

	e.addTrace()
	return e
}

func fromProto(pberr *errpb.ErrorX, serviceName string) *ErrorX {
	details := make(M, len(pberr.GetDetails()))
	for k, v := range pberr.GetDetails() {
		details[fmt.Sprintf("%s.%s", serviceName, k)] = v
	}

	err := &ErrorX{
		Code:    pberr.GetCode(),
		Message: pberr.GetMessage(),
		Type:    Type(pberr.GetType()),
		Details: details,
		origin:  errors.New(pberr.GetMessage()),
		trace:   pberr.GetTrace(),
	}

	return err
}
