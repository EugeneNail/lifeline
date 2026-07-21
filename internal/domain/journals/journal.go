package journals

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/google/uuid"
)

// Journal represents a daily journal entry orchestrating validated value objects.
type Journal struct {
	date      domain.Date
	note      Note
	createdAt time.Time
	updatedAt time.Time
	accountId auth.ID
}

// TODO rename raw constructors of the other domain models
// NewFromRaw returns a validated journal or domain validation violations when construction fails.
func NewFromRaw(rawDate time.Time, rawNote string, accountId auth.ID) (*Journal, domain.Violations) {
	violations := domain.NewViolations()

	date, violation := domain.NewDate(rawDate)
	if violation != nil {
		violations.Add("date", violation)
	}

	note, violation := NewNote(rawNote)
	if violation != nil {
		violations.Add("note", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	now := time.Now()

	return &Journal{
		date:      date,
		note:      note,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// New returns a validated journal.
func New(date domain.Date, note Note, accountId auth.ID) *Journal {
	now := time.Now()

	return &Journal{
		date:      date,
		note:      note,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}
}

// Restore returns a journal reconstructed from persisted values without validating them.
func Restore(date time.Time, note string, createdAt time.Time, updatedAt time.Time, accountId uuid.UUID) *Journal {
	return &Journal{
		date:      domain.Date(date),
		note:      Note(note),
		createdAt: createdAt,
		updatedAt: updatedAt,
		accountId: auth.ID(accountId),
	}
}

// Date returns the journal date.
func (journal *Journal) Date() domain.Date {
	return journal.date
}

// Note returns the journal note.
func (journal *Journal) Note() Note {
	return journal.note
}

// ChangeNote updates the journal note and refreshes the modification timestamp.
func (journal *Journal) ChangeNote(note Note) {
	journal.note = note
	journal.updatedAt = time.Now()
}

// CreatedAt returns the timestamp when the journal was created.
func (journal *Journal) CreatedAt() time.Time {
	return journal.createdAt
}

// UpdatedAt returns the timestamp when the journal was last updated.
func (journal *Journal) UpdatedAt() time.Time {
	return journal.updatedAt
}

// AccountId returns the account identifier that owns the journal.
func (journal *Journal) AccountId() auth.ID {
	return journal.accountId
}
