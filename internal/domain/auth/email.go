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

// NewEmail validates, normalizes, and returns a domain email value.
func NewEmail(rawEmail string) (Email, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(rawEmail))
	if normalizedEmail == "" {
		return "", domain.NewError("email is empty")
	}

	// TODO extract into a constant and add as a placeholder into the error message
	if len(normalizedEmail) > 200 {
		return "", domain.NewError("email length exceeds 200 characters")
	}

	parsedEmail, err := mail.ParseAddress(normalizedEmail)
	if err != nil || parsedEmail.Address != normalizedEmail {
		return "", domain.NewError("email has invalid format")
	}

	return Email(normalizedEmail), nil
}
