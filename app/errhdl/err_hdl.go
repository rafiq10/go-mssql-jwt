package errhdl

import "net/http"

var (
	ErrAuth      = &ApiError{status: http.StatusUnauthorized, msg: "invalid token"}
	ErrNotFound  = &ApiError{status: http.StatusNotFound, msg: "not found"}
	ErrDuplicate = &ApiError{status: http.StatusBadRequest, msg: "duplicate"}
)

// ApiError returns HTTP status code and an error message
type ApiErr interface {
	ApiError() (int, string)
}

// ApiError type represents an error with an associated HTTP status code and message
type ApiError struct {
	status int
	msg    string
}

// Allows StatusError to satisfy the error interface
func (e ApiError) Error() string {
	return e.msg
}

// Returns HTTP status code and message
func (e ApiError) ApiErr() (int, string) {
	return e.status, e.msg
}

type apiWrappedError struct {
	error
	source   string
	apiError *ApiError
}

func (e apiWrappedError) Is(err error) bool {
	return e.apiError == err
}

func (e apiWrappedError) ApiError() (int, string) {
	return e.apiError.ApiErr()
}

func WrapError(err error, apiWrapp *ApiError, src string) error {
	return apiWrappedError{error: err, source: src, apiError: apiWrapp}
}
