package error

import "net/http"

type ResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ResultError) Error() string {
	return e.Message
}

var (
	BadRequest        = &ResultError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrUnauthorized   = &ResultError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden      = &ResultError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrNotFound       = &ResultError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrInternalServer = &ResultError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrDatabase       = &ResultError{Code: http.StatusInternalServerError, Message: "database error"}
)
