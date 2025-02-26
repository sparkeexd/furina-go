package errors

import (
	"fmt"
	"net/http"
)

// Custom error for handling internal errors and HoYoLab endpoint errors.
// Errors from HoYoLab endpoints have special error codes within the response body after returning a 200 OK response.
type Error struct {
	ErrorCode    int
	ErrorMessage string
}

// Constructor.
func NewError(errorCode int, errorMessage string) error {
	err := Error{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}

	return err
}

// Returns custom error message.
// Implements from built-in error interface.
func (err Error) Error() string {
	return fmt.Sprintf("%d %s\n%s", err.ErrorCode, http.StatusText(err.ErrorCode), err.ErrorMessage)
}
