package create_journal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/domain/journals"
)

// Handler executes the create-journal use case.
type Handler struct {
	journals journals.Repository
}

// NewHandler returns a create-journal handler configured with the journal repository or an error when a dependency is missing.
func NewHandler(journals journals.Repository) (*Handler, error) {
	if journals == nil {
		return nil, fmt.Errorf("create_journal handler requires a journal repository")
	}

	return &Handler{journals: journals}, nil
}

// Command carries the data required to create a daily journal.
type Command struct {
	Date      time.Time
	Note      string
	AccountID auth.ID
}

// Handle validates the command, updates an existing journal or creates a new one, and returns an error when persistence fails.
func (handler *Handler) Handle(ctx context.Context, command Command) error {
	violations := domain.NewViolations()

	date, err := domain.NewDate(command.Date)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return fmt.Errorf("creating a journal date: %w", err)
		}

		violations.Add("date", violation)
	}

	note, err := journals.NewNote(command.Note)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return fmt.Errorf("creating a journal note: %w", err)
		}

		violations.Add("note", violation)
	}

	if violations.HasViolations() {
		return violations
	}

	journal, err := handler.journals.Find(ctx, journals.NewFilter().
		WithAccountIds(command.AccountID).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return fmt.Errorf("finding a journal by account id %q and date %q: %w", command.AccountID, time.Time(date), err)
	}

	isNew := false
	if journal == nil {
		isNew = true
		journal = journals.New(date, note, command.AccountID)
	}

	journal.ChangeNote(note)

	if isNew {
		if err := handler.journals.Add(ctx, journal); err != nil {
			return fmt.Errorf("adding a new journal to the collection: %w", err)
		}
	} else {
		if err := handler.journals.Update(ctx, journal); err != nil {
			return fmt.Errorf("updating a journal by account id %q and date %q: %w", command.AccountID, time.Time(date), err)
		}
	}

	return nil
}
