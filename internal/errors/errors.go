package errors

import (
	"fmt"
	"log"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    // HTTP status code
	Message string // User-facing message
	Details string // Internal details for logging
	Err     error  // Original error
}

func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string, details string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

// HandleError handles errors and sends appropriate HTTP response
func HandleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *AppError:
		// Log internal details
		if e.Details != "" {
			log.Printf("Error [%d]: %s - Details: %s", e.Code, e.Message, e.Details)
		} else if e.Err != nil {
			log.Printf("Error [%d]: %s - %v", e.Code, e.Message, e.Err)
		}
		// Send user-facing message only
		http.Error(w, e.Message, e.Code)
	default:
		// Unknown error - log details but don't expose to client
		log.Printf("Internal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// WrapError wraps an error with context
func WrapError(err error, message string, code int) *AppError {
	return NewAppError(code, message, err.Error(), err)
}

// InternalError creates an internal server error
func InternalError(details string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "Internal server error", details, err)
}

// NotFoundError creates a not found error
func NotFoundError(resource string) *AppError {
	return NewAppError(http.StatusNotFound, fmt.Sprintf("%s not found", resource), "", nil)
}

// BadRequestError creates a bad request error
func BadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, "", nil)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return NewAppError(http.StatusUnauthorized, message, "", nil)
}

