package uf

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// get the error in string format
func (e HttpError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Code, http.StatusText(e.Code), e.Message)
}

// 400 Bad Request error
func BadRequest(m string) HttpError {
	return HttpError{http.StatusBadRequest, m}
}

// 401 Unauthorized error
func Unauthorized(m string) HttpError {
	return HttpError{http.StatusUnauthorized, m}
}

// 403 Forbidden error
func Forbidden(m string) HttpError {
	return HttpError{http.StatusForbidden, m}
}

// 404 Not Found error
func NotFound(m string) HttpError {
	return HttpError{http.StatusNotFound, m}
}

// 405 Method Not Allowed error
func MethodNotAllowed(m string) HttpError {
	return HttpError{http.StatusMethodNotAllowed, m}
}

// 500 Internal Server Error error
func InternalServerError(m string) HttpError {
	return HttpError{http.StatusInternalServerError, m}
}
