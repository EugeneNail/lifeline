package journal

import (
	"errors"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Journal struct {
	id        ID
	date      Date
	mood      Mood
	note      Note
	createdAt time.Time
	updatedAt time.Time
	accountId auth.ID
}

func NewJournal(rawDate time.Time, rawMood int, rawNote string, accountId auth.ID) (*Journal, error) {
	violations := domain.NewViolations()

	date, err := NewDate(rawDate)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return nil, fmt.Errorf("creating a date: %w", err)
		}

		violations.Add("date", violation)
	}

	mood, err := NewMood(rawMood)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return nil, fmt.Errorf("creating a mood: %w", err)
		}

		violations.Add("mood", violation)
	}

	note, err := NewNote(rawNote)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return nil, fmt.Errorf("creating a note: %w", err)
		}

		violations.Add("note", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	now := time.Now()

	return &Journal{
		id:        NewID(),
		date:      date,
		mood:      mood,
		note:      note,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// RestoreJournal returns a journal reconstructed from persisted values without validating or changing them.
func RestoreJournal(id uuid.UUID, date time.Time, mood int, note string, createdAt time.Time, updatedAt time.Time, accountId uuid.UUID) *Journal {
	return &Journal{
		id:        ID(id),
		date:      Date(date),
		mood:      Mood(mood),
		note:      Note(note),
		createdAt: createdAt,
		updatedAt: updatedAt,
		accountId: auth.ID(accountId),
	}
}

// ===================================== ID ============================================
// =====================================================================================

// TODO remove as a custom type
type ID uuid.UUID

// NilID is the zero-value user identifier.
var NilID = ID(uuid.Nil)

// ID returns the user identifier.
func (journal *Journal) ID() ID {
	return journal.id
}

// Uuid returns the UUID value of the identifier.
func (id ID) Uuid() uuid.UUID {
	return uuid.UUID(id)
}

// NewID returns a new UUIDv7-based user identifier.
func NewID() ID {
	return ID(uuidv7.New())
}

// String returns the UUID string representation of the identifier.
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// ===================================== Date ============================================
// =======================================================================================

type Date time.Time

func NewDate(raw time.Time) (Date, error) {
	if raw.IsZero() {
		return Date{}, domain.NewViolation("date is empty")
	}

	date := raw.Truncate(time.Hour * 24)
	minDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, date.Location())
	maxDate := time.Date(2099, time.January, 1, 0, 0, 0, 0, date.Location())

	if date.Before(minDate) || date.After(maxDate) {
		return Date{}, domain.NewViolation("date must be between 2000-01-01 and 2099-01-01")
	}

	return Date(date), nil
}

func (journal *Journal) Date() Date {
	return journal.date
}

// ===================================== Mood ============================================
// =======================================================================================

func (journal *Journal) Mood() Mood {
	return journal.mood
}

func (journal *Journal) ChangeMood(mood Mood) {
	journal.mood = mood
	journal.updatedAt = time.Now()
}

func (mood Mood) String() string {
	return moodLabels[mood]
}

// ===================================== Note ============================================
// =======================================================================================

const noteMaxLength = 10000

type Note string

// NewNote returns a note value built from the provided raw text and an error when note creation fails.
func NewNote(raw string) (Note, error) {
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

func (journal *Journal) Note() Note {
	return journal.note
}

func (journal *Journal) ChangeNote(note Note) {
	journal.note = note
	journal.updatedAt = time.Now()
}

// ===================================== Misc ============================================
// =======================================================================================

func (journal *Journal) CreatedAt() time.Time {
	return journal.createdAt
}

func (journal *Journal) UpdatedAt() time.Time {
	return journal.updatedAt
}

func (journal *Journal) AccountId() auth.ID {
	return journal.accountId
}
