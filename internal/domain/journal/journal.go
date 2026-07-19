package journal

import (
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
)

// Journal represents a daily journal entry orchestrating validated value objects.
type Journal struct {
	id        uuid.UUID
	date      domain.Date
	note      Note
	createdAt time.Time
	updatedAt time.Time
	accountId auth.ID
}

// New returns a validated journal or domain validation violations when construction fails.
func New(rawDate time.Time, rawNote string, accountId auth.ID) (*Journal, error) {
	violations := domain.NewViolations()

	date, err := domain.NewDate(rawDate)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return nil, fmt.Errorf("creating a date: %w", err)
		}

		violations.Add("date", violation)
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
		id:        uuid.UUID(uuidv7.New()),
		date:      date,
		note:      note,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// Restore returns a journal reconstructed from persisted values without validating them.
func Restore(id uuid.UUID, date time.Time, note string, createdAt time.Time, updatedAt time.Time, accountId uuid.UUID) *Journal {
	return &Journal{
		id:        id,
		date:      domain.Date(date),
		note:      Note(note),
		createdAt: createdAt,
		updatedAt: updatedAt,
		accountId: auth.ID(accountId),
	}
}

// ID returns the journal identifier.
func (journal *Journal) ID() uuid.UUID {
	return journal.id
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
