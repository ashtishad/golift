package common

import (
	"fmt"
	"net/http"
)

type AppError interface {
	Error() string
	Code() int
	Cause(err error) error
}

type Error struct {
	Message    string
	StatusCode int
	Err        error // Field for the wrapped error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Code() int {
	return e.StatusCode
}

func (e *Error) Cause(err error) error {
	if err != nil {
		e.Err = fmt.Errorf("%w", err)
	}
	return e
}

func NewBadRequestError(message string) AppError {
	return &Error{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewInternalServerError(message string, err error) AppError {
	return &Error{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err, // Wrap the internal error
	}
}

// NewNotFoundError creates a new APIError for not found errors.
func NewNotFoundError(message string) AppError {
	return &Error{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError creates a new APIError for unauthorized requests.
func NewUnauthorizedError(message string) AppError {
	return &Error{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewConflictError creates a new APIError for conflict errors.
func NewConflictError(message string) AppError {
	return &Error{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}
