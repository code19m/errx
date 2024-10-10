package errxto

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/code19m/errx"
)

// errto.HTTP is intended for use in HTTP handlers to convert ErrorX instances to HTTP responses.
// If the error is nil, no response is written.
// If the error code is INTERNAL it won't be written error details in response.
func HTTP(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if _, ok := err.(*errx.ErrorX); !ok {
		err = errx.Wrap(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode(err))
	writeBody(w, err)
}

func writeBody(w http.ResponseWriter, err error) {
	e, ok := err.(*errx.ErrorX)
	if ok && e != nil && e.Code != errx.CodeInternal {
		jsonBody, marshalErr := json.Marshal(e)
		if marshalErr != nil {
			// If we can't marshal the error, we should panic because it's a bug in the code.
			// And at this point, we can't do anything else.
			panic(fmt.Errorf("marshal error: %v. original error: %v", marshalErr, err))
		}
		_, _ = w.Write(jsonBody)
		return
	}
	_, _ = w.Write(
		[]byte(fmt.Sprintf(`{"code": "%s", "message": "%s"}`, errx.CodeInternal, "Internal server error")),
	)
}

func httpStatusCode(err error) int {
	if e, ok := err.(*errx.ErrorX); ok {
		switch e.Type {
		case errx.T_Validation:
			return http.StatusBadRequest
		case errx.T_NotFound:
			return http.StatusNotFound
		case errx.T_Conflict:
			return http.StatusConflict
		case errx.T_Internal:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}
