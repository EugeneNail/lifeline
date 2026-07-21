package auth

import (
	"net/mail"
	"strings"

	"github.com/EugeneNail/lifeline/internal/domain"
)

// Email represents a validated user email address.
type Email string

// Email returns the user's email address.
func (account *Account) Email() Email {
	return account.email
}

// String returns the string representation of the email address.
func (email Email) String() string {
	return string(email)
}

// ChangeEmail updates the user's email address.
func (account *Account) ChangeEmail(email Email) {
	account.email = email
}

// NewEmail validates, normalizes, and returns a domain email value or a violation when the email is invalid.
func NewEmail(rawEmail string) (Email, domain.Violation) {
	email := strings.ToLower(strings.TrimSpace(rawEmail))
	if email == "" {
		return "", domain.NewViolation("email is empty")
	}

	// TODO extract into a constant and add as a placeholder into the error message
	if len(email) > 200 {
		return "", domain.NewViolation("email length exceeds 200 characters")
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return "", domain.NewViolation("email has invalid format")
	}

	return Email(email), nil
}
