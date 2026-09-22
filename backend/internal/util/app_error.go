package util

import "fmt"

// AppError is a business error with HTTP status and business code.
type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

// NewAppError builds an AppError.
func NewAppError(httpStatus, code int, message string) *AppError {
	return &AppError{HTTPStatus: httpStatus, Code: code, Message: message}
}
