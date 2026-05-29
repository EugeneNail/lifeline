package auth

import (
	"net/mail"
	"strings"

	"github.com/EugeneNail/lifeline/internal/domain"
)

// Email represents a validated user email address.
type Email string

// Email returns the user's email address.
func (user *Account) Email() Email {
	return user.email
}

// ChangeEmail updates the user's email address.
func (user *Account) ChangeEmail(email Email) {
	user.email = email
}

// NewEmail validates, normalizes, and returns a domain email value.
func NewEmail(rawEmail string) (Email, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(rawEmail))
	if normalizedEmail == "" {
		return "", domain.NewError("email is empty")
	}

	if len(normalizedEmail) > 200 {
		return "", domain.NewError("email length exceeds 200 characters")
	}

	parsedEmail, err := mail.ParseAddress(normalizedEmail)
	if err != nil || parsedEmail.Address != normalizedEmail {
		return "", domain.NewError("email has invalid format")
	}

	return Email(normalizedEmail), nil
}
