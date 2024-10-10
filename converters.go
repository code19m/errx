package errx

import (
	"errors"

	"github.com/code19m/errx/internal/errpb"
	"github.com/jackc/pgconn"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx/v5"
)

func FromPgxQuery(err error, notFoundCode string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		e := New(T_NotFound, notFoundCode, "Resource not found. Code: "+notFoundCode)
		e.addTrace()
		return e
	}

	e := New(T_Internal, CodeInternal, err.Error())
	e.addTrace()
	return e
}

func FromPgxExec(err error, conflictCode string) error {
	if err == nil {
		return nil
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		e := New(T_Internal, CodeInternal, err.Error())
		e.addTrace()
		return e
	}

	var e *ErrorX
	switch pgErr.Code {
	case "23505":
		e = New(T_Conflict, conflictCode, "Resource already exists. Code: "+conflictCode).
			WithDetails(M{"pg_constraint": pgErr.ConstraintName})
	default:
		e = New(T_Internal, CodeInternal, err.Error()).
			WithDetails(M{
				"pg_hint":       pgErr.Hint,
				"pg_detail":     pgErr.Detail,
				"pg_constraint": pgErr.ConstraintName,
				"pg_column":     pgErr.ColumnName,
			})
	}

	e.addTrace()
	return e
}

func FromGRPC(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		e := New(T_Internal, CodeInternal, err.Error())
		e.addTrace()
		return e
	}

	for _, detail := range st.Details() {
		if pb, ok := detail.(*errpb.ErrorX); ok {
			e := fromProto(pb)
			e.addTrace()
			return e
		}
	}
	e := New(T_Internal, CodeInternal, err.Error()).WithDetails(M{
		"grpc_code": st.Code().String(),
		"grpc_msg":  st.Message(),
	})
	e.addTrace()
	return e
}

func fromProto(pberr *errpb.ErrorX) *ErrorX {
	err := &ErrorX{
		Code:    pberr.GetCode(),
		Message: pberr.GetMessage(),
		Type:    Type(pberr.GetType()),
		Details: pberr.GetDetails(),
		origin:  errors.New(pberr.GetMessage()),
		trace:   pberr.GetTrace(),
	}
	return err
}
