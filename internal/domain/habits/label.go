package habits

import (
	"strings"
	"unicode/utf8"

	"github.com/EugeneNail/lifeline/internal/domain"
)

const (
	habitLabelMinLength = 3
	habitLabelMaxLength = 32
)

// NewLabel returns a normalized habit label or a domain error when the label violates domain rules.
func NewLabel(rawLabel string) (string, error) {
	label := strings.TrimSpace(rawLabel)
	length := utf8.RuneCountInString(label)

	if length < habitLabelMinLength || length > habitLabelMaxLength {
		return "", domain.NewErrorf(
			"label length must be between %d and %d characters",
			habitLabelMinLength,
			habitLabelMaxLength,
		)
	}

	for _, character := range label {
		if IsAllowedLabelCharacter(character) {
			continue
		}

		return "", domain.NewError("label contains unsupported characters")
	}

	return label, nil
}

// IsAllowedLabelCharacter reports whether the character is allowed in a habit label.
func IsAllowedLabelCharacter(character rune) bool {
	if character >= 'a' && character <= 'z' {
		return true
	}

	if character >= 'A' && character <= 'Z' {
		return true
	}

	if character >= '0' && character <= '9' {
		return true
	}

	switch character {
	case '-', ' ', ',', '!', '?':
		return true
	default:
		return false
	}
}
