package application

import (
	"errors"
	"github.com/EugeneNail/lifeline/internal/domain"
)

// FieldErrors collects field-level validation messages keyed by field name.
type FieldErrors struct {
	errors map[string]string
}

// NewFieldErrors returns an empty collection of field validation errors.
func NewFieldErrors() FieldErrors {
	return FieldErrors{
		make(map[string]string),
	}
}

// Error returns a generic validation failure message.
func (fe FieldErrors) Error() string {
	return "field validation failed"
}

// AddFromDomain stores a domain validation error under the provided field and returns the original error when it is not a domain error.
func (fe FieldErrors) AddFromDomain(field string, err error) error {
	var domainError domain.Error
	if errors.As(err, &domainError) {
		fe.Add(field, domainError.Error())
		return nil
	}

	return err
}

// Add stores a validation message for the provided field.
func (fe FieldErrors) Add(field string, message string) {
	fe.errors[field] = message
}

// HasErrors reports whether at least one field validation error has been collected.
func (fe FieldErrors) HasErrors() bool {
	return len(fe.errors) > 0
}

// Errors returns the collected field validation errors keyed by field name.
func (fe FieldErrors) Errors() map[string]string {
	return fe.errors
}
