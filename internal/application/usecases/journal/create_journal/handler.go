package create_journal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/domain/journal"
	"github.com/google/uuid"
)

// Handler executes the create-journal use case.
type Handler struct {
	journals              journal.JournalRepository
	journalCreationPolicy *journal.CreationPolicy
}

// NewHandler returns a create-journal handler configured with the journal repository and creation policy or an error when a dependency is missing.
func NewHandler(journals journal.JournalRepository, journalCreationPolicy *journal.CreationPolicy) (*Handler, error) {
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
	Note      string
	AccountID auth.ID
}

// Handle validates the command, stores a new daily journal, and returns the journal identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	journalEntry, err := journal.New(command.Date, command.Note, command.AccountID)
	if err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return uuid.Nil, violations
		}

		return uuid.Nil, fmt.Errorf("creating a journal: %w", err)
	}

	if err := handler.journalCreationPolicy.Check(ctx, command.AccountID, command.Date); err != nil {
		if errors.Is(err, journal.ErrDateIsOccupied) {
			return uuid.Nil, err
		}

		return uuid.Nil, fmt.Errorf("ensuring journal can be added: %w", err)
	}

	if err := handler.journals.Add(ctx, journalEntry); err != nil {
		return uuid.Nil, fmt.Errorf("adding a new journal to the collection: %w", err)
	}

	return journalEntry.ID(), nil
}
