package errors

import (
	"fmt"
	"invoice-scan/backend/pkg"
	"invoice-scan/backend/pkg/ulid"
	"net/http"
)

var (
	InvalidDataErr = func(code, message string) *IError {
		return NewHTTPErrorCode(http.StatusBadRequest, code, message)
	}
	NotfoundErr = func(code, message string) *IError {
		return NewHTTPErrorCode(http.StatusNotFound, code, message)
	}
	ForbiddenErr = func(code, message string) *IError {
		return NewHTTPErrorCode(http.StatusForbidden, code, message)
	}
	UnauthorizedErr = func(code, message string) *IError {
		return NewHTTPErrorCode(http.StatusUnauthorized, code, message)
	}
	InternalServerErr = func(code, message string) *IError {
		return NewHTTPErrorCode(http.StatusInternalServerError, code, message)
	}
)

type IError struct {
	ErrorCode
	httpCode int
	id       string
	data     interface{}
}

func (e *IError) ID() string {
	return e.id
}

func (e *IError) HTTPCode() int {
	if e.httpCode <= 0 {
		return http.StatusBadRequest
	}
	return e.httpCode
}

func (e *IError) Data() interface{} {
	return e.data
}

func (e *IError) WithData(d interface{}) *IError {
	e.data = d
	return e
}

func (e *IError) JSON() map[string]interface{} {
	return map[string]interface{}{
		"error_code":    e.defaultErrorCode(),
		"error_summary": e.message,
		"error_id":      e.id,
	}
}

func (e *IError) defaultErrorCode() string {
	if pkg.IsStringEmpty(e.code) {
		return fmt.Sprintf("E%d", e.httpCode)
	}
	return e.code
}

// NewHTTPErrorCode create error with http code and message
func NewHTTPErrorCode(httpCode int, code, message string) *IError {
	return &IError{
		ErrorCode: ErrorCode{
			code:    code,
			message: message,
		},
		httpCode: httpCode,
		id:       "err." + ulid.GenerateULID(),
	}
}
