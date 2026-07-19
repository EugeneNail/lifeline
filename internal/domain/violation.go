package domain

import "fmt"

// Violation represents a domain-level validation or business-rule violation.
type Violation interface {
	error
}

type violation struct {
	message string
}

// Error returns the domain failure message.
func (violation violation) Error() string {
	return violation.message
}

// NewViolation returns a domain error with the provided message.
func NewViolation(message string) Violation {
	return violation{message: message}
}

// NewViolationf returns a domain error with a formatted message.
func NewViolationf(format string, a ...any) Violation {
	return violation{message: fmt.Sprintf(format, a...)}
}
