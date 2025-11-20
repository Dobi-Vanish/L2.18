package errors

import "fmt"

// ValidationError 400 error.
type ValidationError struct {
	Field   string
	Message string
}

// // Error to provide 400 error messages.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// BusinessError 503 error.
type BusinessError struct {
	Operation string
	Message   string
}

// Error to provide 503 error messages.
func (e BusinessError) Error() string {
	return fmt.Sprintf("business error: %s - %s", e.Operation, e.Message)
}

// InternalError 500 error.
type InternalError struct {
	Operation string
	Message   string
}

// Error to provide 500 error messages.
func (e InternalError) Error() string {
	return fmt.Sprintf("internal error: %s - %s", e.Operation, e.Message)
}
