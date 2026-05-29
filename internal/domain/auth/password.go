package auth

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/EugeneNail/lifeline/internal/domain"
)

// Password represents a validated raw user password.
type Password string

// Password returns the account password hash.
func (account *Account) Password() HashedPassword {
	return account.password
}

// NewPassword validates and returns a password value when it satisfies the domain password policy.
func NewPassword(rawPassword string) (Password, error) {
	if utf8.RuneCountInString(rawPassword) < 8 {
		return "", domain.NewError("password must be at least 8 characters long")
	}

	if utf8.RuneCountInString(rawPassword) > 128 {
		return "", domain.NewError("password must be at most 128 characters long")
	}

	var hasUppercase bool
	var hasLowercase bool
	var hasDigit bool
	var hasSpecial bool

	for _, symbol := range rawPassword {
		switch {
		case unicode.IsUpper(symbol):
			hasUppercase = true
		case unicode.IsLower(symbol):
			hasLowercase = true
		case unicode.IsDigit(symbol):
			hasDigit = true
		case unicode.IsPunct(symbol), unicode.IsSymbol(symbol):
			hasSpecial = true
		}
	}

	if !hasUppercase || !hasLowercase || !hasDigit || !hasSpecial {
		return "", domain.NewError(passwordPolicyMessage(hasUppercase, hasLowercase, hasDigit, hasSpecial))
	}

	return Password(rawPassword), nil
}

func passwordPolicyMessage(hasUppercase, hasLowercase, hasDigit, hasSpecial bool) string {
	missingRequirements := make([]string, 0, 4)

	if !hasUppercase {
		missingRequirements = append(missingRequirements, "uppercase letter")
	}

	if !hasLowercase {
		missingRequirements = append(missingRequirements, "lowercase letter")
	}

	if !hasDigit {
		missingRequirements = append(missingRequirements, "digit")
	}

	if !hasSpecial {
		missingRequirements = append(missingRequirements, "special symbol")
	}

	return fmt.Sprintf("password must contain at least one %s", joinRequirements(missingRequirements))
}

func joinRequirements(requirements []string) string {
	switch len(requirements) {
	case 0:
		return ""
	case 1:
		return requirements[0]
	case 2:
		return requirements[0] + " and " + requirements[1]
	default:
		result := requirements[0]
		for index := 1; index < len(requirements)-1; index++ {
			result += ", " + requirements[index]
		}

		return result + ", and " + requirements[len(requirements)-1]
	}
}

// Hash returns the hashed password produced by the provided PasswordHasher.
func (password Password) Hash(hasher PasswordHasher) (HashedPassword, error) {
	if hasher == nil {
		return "", fmt.Errorf("password hasher is nil")
	}

	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	return hashedPassword, nil
}
