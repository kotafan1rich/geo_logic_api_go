package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

var (
	ErrNotFound     = &AppError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: "Resource not found"}
	ErrBadRequest   = &AppError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: "Invalid request"}
	ErrInternal     = &AppError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Internal server error"}
	ErrUnauthorized = &AppError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrForbidden    = &AppError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: "Forbidden"}
	ErrConflict     = &AppError{Status: http.StatusConflict, Code: "CONFLICT", Message: "Conflict"}
)

func New(status int, code, message, details string) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Details: details}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func Wrap(err error, appErr *AppError) *AppError {
	if err == nil {
		return appErr
	}
	return &AppError{
		Status:  appErr.Status,
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: err.Error(),
		Err:     err,
	}
}

func ValidationError(details string) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: details,
	}
}
