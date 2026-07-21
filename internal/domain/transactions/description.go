package transactions

import (
	"unicode/utf8"

	"github.com/EugeneNail/lifeline/internal/domain"
)

const (
	descriptionMinLength = 1
	descriptionMaxLength = 32
)

// Description is a validated transaction description.
type Description string

// NewDescription returns a validated transaction description or a violation when the description is invalid.
func NewDescription(rawDescription string) (Description, domain.Violation) {
	length := utf8.RuneCountInString(rawDescription)
	if length < descriptionMinLength || length > descriptionMaxLength {
		return "", domain.NewViolationf(
			"description length must be between %d and %d characters",
			descriptionMinLength,
			descriptionMaxLength,
		)
	}

	for _, character := range rawDescription {
		if isAllowedDescriptionCharacter(character) {
			continue
		}

		return "", domain.NewViolation("description contains unsupported characters")
	}

	return Description(rawDescription), nil
}

// String returns the raw transaction description.
func (description Description) String() string {
	return string(description)
}

func isAllowedDescriptionCharacter(character rune) bool {
	if character >= '0' && character <= '9' {
		return true
	}

	if character >= 'a' && character <= 'z' {
		return true
	}

	if character >= 'A' && character <= 'Z' {
		return true
	}

	switch character {
	case ' ', '.', ',', '\'', '"':
		return true
	default:
		return false
	}
}
