package domain

import "fmt"

// Error represents a domain-level validation or business-rule violation.
type Error interface {
	error
}

type domainError struct {
	message string
}

// Error returns the domain failure message.
func (error domainError) Error() string {
	return error.message
}

// NewError returns a domain error with the provided message.
func NewError(message string) Error {
	return domainError{message: message}
}

// NewErrorf returns a domain error with a formatted message.
func NewErrorf(format string, a ...any) Error {
	return domainError{message: fmt.Sprintf(format, a...)}
}
