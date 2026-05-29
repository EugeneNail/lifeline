package domain

// Error represents a domain-level validation or business-rule violation.
type Error struct {
	message string
}

// Error returns the domain failure message.
func (error Error) Error() string {
	return error.message
}

// NewError returns a domain error with the provided message.
func NewError(message string) Error {
	return Error{message: message}
}
