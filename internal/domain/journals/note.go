package journals

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/EugeneNail/lifeline/internal/domain"
)

const noteMaxLength = 10000

// Note is a validated free-form journal note.
type Note string

// NewNote returns a validated note or a violation when the note text is invalid.
func NewNote(raw string) (Note, domain.Violation) {
	raw = normalizeNote(raw)

	if !utf8.ValidString(raw) {
		return "", domain.NewViolation("note must be valid UTF-8")
	}

	// Emojis may look like one symbol while occupying 2-4 runes, so this check counts runes rather than visual glyphs.
	if utf8.RuneCountInString(raw) > noteMaxLength {
		return "", domain.NewViolationf("note must not exceed %d characters", noteMaxLength)
	}

	if containsForbiddenControlChars(raw) {
		return "", domain.NewViolation("note contains forbidden control characters")
	}

	return Note(raw), nil
}

func normalizeNote(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	return text
}

// containsForbiddenControlChars reports whether the note text contains control characters other than newline and tab.
func containsForbiddenControlChars(text string) bool {
	for _, character := range text {
		if character == '\n' || character == '\t' {
			continue
		}

		if unicode.IsControl(character) {
			return true
		}
	}

	return false
}
