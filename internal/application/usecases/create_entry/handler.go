package create_entry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/domain/entries"
)

// Handler executes the create-entry use case.
type Handler struct {
	entries             entries.EntryRepository
	entryCreationPolicy *entries.EntryCreationPolicy
}

// NewHandler returns a create-entry handler configured with the entry repository and creation policy or an error when a dependency is missing.
func NewHandler(entries entries.EntryRepository, entryCreationPolicy *entries.EntryCreationPolicy) (*Handler, error) {
	if entries == nil {
		return nil, fmt.Errorf("create_entry handler requires an entry repository")
	}

	if entryCreationPolicy == nil {
		return nil, fmt.Errorf("create_entry handler requires an entry creation policy")
	}

	return &Handler{
		entries:             entries,
		entryCreationPolicy: entryCreationPolicy,
	}, nil
}

// Command carries the data required to create a daily entry.
type Command struct {
	Date      time.Time
	Mood      int
	Note      string
	AccountID auth.ID
}

// Handle validates the command, stores a new daily entry, and returns the entry identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (entries.ID, error) {
	entry, err := entries.New(command.Date, command.Mood, command.Note, command.AccountID)
	if err != nil {
		var validationErrors domain.ValidationErrors
		if errors.As(err, &validationErrors) {
			return entries.NilID, validationErrors
		}

		return entries.NilID, fmt.Errorf("creating an entry: %w", err)
	}

	if err := handler.entryCreationPolicy.EnsureCanAddEntry(ctx, command.AccountID, command.Date); err != nil {
		if errors.Is(err, entries.ErrDateIsOccupied) {
			return entries.NilID, err
		}

		return entries.NilID, fmt.Errorf("ensuring entry can be added: %w", err)
	}

	if err := handler.entries.Add(ctx, entry); err != nil {
		return entries.NilID, fmt.Errorf("adding a new entry to the collection: %w", err)
	}

	return entry.ID(), nil
}
