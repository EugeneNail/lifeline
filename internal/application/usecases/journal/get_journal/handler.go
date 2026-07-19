package get_journal

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/domain/journals"
)

// Handler executes the get-journal use case.
type Handler struct {
	journals journals.Repository
}

// NewHandler returns a get-journal handler configured with the journal repository or an error when the dependency is missing.
func NewHandler(journalsRepository journals.Repository) (*Handler, error) {
	if journalsRepository == nil {
		return nil, fmt.Errorf("get_journal handler requires a journal repository")
	}

	return &Handler{journals: journalsRepository}, nil
}

// Query carries the data required to load a journal for a specific day.
type Query struct {
	AccountID auth.ID
	Date      time.Time
}

// Handle returns the journal matching the query or nil when no journal exists, or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (*journals.Journal, error) {
	journalEntry, err := handler.journals.Find(ctx, journals.NewFilter().
		WithAccountIds(query.AccountID).
		WithDates(query.Date),
	)
	if err != nil {
		return nil, fmt.Errorf("finding a journal by account id %q and date %q: %w", query.AccountID, query.Date, err)
	}

	return journalEntry, nil
}
