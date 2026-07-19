package create_journal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/domain/journal"
)

// Handler executes the create-journal use case.
type Handler struct {
	journals              journal.JournalRepository
	journalCreationPolicy *journal.JournalCreationPolicy
}

// NewHandler returns a create-journal handler configured with the journal repository and creation policy or an error when a dependency is missing.
func NewHandler(journals journal.JournalRepository, journalCreationPolicy *journal.JournalCreationPolicy) (*Handler, error) {
	if journals == nil {
		return nil, fmt.Errorf("create_journal handler requires a journal repository")
	}

	if journalCreationPolicy == nil {
		return nil, fmt.Errorf("create_journal handler requires a journal creation policy")
	}

	return &Handler{
		journals:              journals,
		journalCreationPolicy: journalCreationPolicy,
	}, nil
}

// Command carries the data required to create a daily journal.
type Command struct {
	Date      time.Time
	Mood      int
	Note      string
	AccountID auth.ID
}

// Handle validates the command, stores a new daily journal, and returns the journal identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (journal.ID, error) {
	journalEntry, err := journal.NewJournal(command.Date, command.Mood, command.Note, command.AccountID)
	if err != nil {
		var validationErrors domain.ValidationErrors
		if errors.As(err, &validationErrors) {
			return journal.NilID, validationErrors
		}

		return journal.NilID, fmt.Errorf("creating a journal: %w", err)
	}

	if err := handler.journalCreationPolicy.EnsureCanAdd(ctx, command.AccountID, command.Date); err != nil {
		if errors.Is(err, journal.ErrDateIsOccupied) {
			return journal.NilID, err
		}

		return journal.NilID, fmt.Errorf("ensuring journal can be added: %w", err)
	}

	if err := handler.journals.Add(ctx, journalEntry); err != nil {
		return journal.NilID, fmt.Errorf("adding a new journal to the collection: %w", err)
	}

	return journalEntry.ID(), nil
}
