package newerrors

import "fmt"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// Wrap adds context to the original error
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
