package service

import "fmt"

// ValidationError represents an error validating requestParameters. This
// specific error type should be used to allow the broker's framework to
// differentiate between validation errors and other common, unexpected errors.
type ValidationError struct {
	Field string
	Issue string
}

// NewValidationError returns a new ValidationError for the given field and
// issue
func NewValidationError(field, issue string) *ValidationError {
	return &ValidationError{
		Field: field,
		Issue: issue,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Error validating field '%s': %s", e.Field, e.Issue)
}
