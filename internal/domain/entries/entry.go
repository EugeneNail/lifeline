package entries

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

type Entry struct {
	id        ID
	date      Date
	mood      Mood
	note      Note
	createdAt time.Time
	updatedAt time.Time
	accountId auth.ID
}

func New(rawDate time.Time, rawMood int, rawNote string, accountId auth.ID) (*Entry, error) {
	errs := domain.NewValidationErrors()

	date, err := NewDate(rawDate)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a date: %w", err)
		}

		errs.Add("date", domainError)
	}

	mood, err := NewMood(rawMood)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a mood: %w", err)
		}

		errs.Add("mood", domainError)
	}

	note, err := NewNote(rawNote)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a note: %w", err)
		}

		errs.Add("note", domainError)
	}

	if errs.HasErrors() {
		return nil, errs
	}

	now := time.Now()

	return &Entry{
		id:        NewID(),
		date:      date,
		mood:      mood,
		note:      note,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// Restore returns an entry reconstructed from persisted values without validating or changing them.
func Restore(id uuid.UUID, date time.Time, mood int, note string, createdAt time.Time, updatedAt time.Time, accountId uuid.UUID) *Entry {
	return &Entry{
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

type ID uuid.UUID

// NilID is the zero-value user identifier.
var NilID = ID(uuid.Nil)

// ID returns the user identifier.
func (entry *Entry) ID() ID {
	return entry.id
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
		return Date{}, domain.NewError("date is empty")
	}

	date := raw.Truncate(time.Hour * 24)
	minDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, date.Location())
	maxDate := time.Date(2099, time.January, 1, 0, 0, 0, 0, date.Location())

	if date.Before(minDate) || date.After(maxDate) {
		return Date{}, domain.NewError("date must be between 2000-01-01 and 2099-01-01")
	}

	return Date(date), nil
}

func (entry *Entry) Date() Date {
	return entry.date
}

// ===================================== Mood ============================================
// =======================================================================================

type Mood int

const (
	MoodAwful Mood = 1
	MoodBad   Mood = 2
	MoodOkay  Mood = 3
	MoodGood  Mood = 4
	MoodGreat Mood = 5
)

var moodLabels = map[Mood]string{
	MoodAwful: "Awful",
	MoodBad:   "Bad",
	MoodOkay:  "Okay",
	MoodGood:  "Good",
	MoodGreat: "Great",
}

func NewMood(rawMood int) (Mood, error) {
	if rawMood < int(MoodAwful) || rawMood > int(MoodGreat) {
		return 0, domain.NewErrorf("mood must be in range between %d and %d", MoodAwful, MoodGreat)
	}

	return Mood(rawMood), nil
}

func (entry *Entry) Mood() Mood {
	return entry.mood
}

func (entry *Entry) ChangeMood(mood Mood) {
	entry.mood = mood
	entry.updatedAt = time.Now()
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
		return "", domain.NewError("note must be valid UTF-8")
	}

	// Emojis may look like one symbol while occupying 2-4 runes, so this check counts runes rather than visual glyphs.
	if utf8.RuneCountInString(raw) > noteMaxLength {
		return "", domain.NewErrorf("note must not exceed %d characters", noteMaxLength)
	}

	if containsForbiddenControlChars(raw) {
		return "", domain.NewError("note contains forbidden control characters")
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

func (entry *Entry) Note() Note {
	return entry.note
}

func (entry *Entry) ChangeNote(note Note) {
	entry.note = note
	entry.updatedAt = time.Now()
}

// ===================================== Misc ============================================
// =======================================================================================

func (entry *Entry) CreatedAt() time.Time {
	return entry.createdAt
}

func (entry *Entry) UpdatedAt() time.Time {
	return entry.updatedAt
}

func (entry *Entry) AccountId() auth.ID {
	return entry.accountId
}
