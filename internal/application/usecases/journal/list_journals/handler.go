package list_journals

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/journals"
)

// Handler executes the list-journals use case.
type Handler struct {
	journals journals.Repository
}

// NewHandler returns a list-journals handler configured with the journal repository or an error when the dependency is missing.
func NewHandler(journalsRepository journals.Repository) (*Handler, error) {
	if journalsRepository == nil {
		return nil, fmt.Errorf("list_journals handler requires a journal repository")
	}

	return &Handler{journals: journalsRepository}, nil
}

// toDate returns the provided time normalized to YYYY-MM-DD.
func toDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

// Handle loads journals for the provided account and date range and returns them or an error when lookup fails, including ErrInvalidDateRange when the from date is after the to date.
func (handler *Handler) Handle(ctx context.Context, query Query) ([]*journals.Journal, error) {
	if query.From.After(query.To) {
		return nil, ErrInvalidDateRange
	}

	dates := make([]time.Time, 0)
	for currentDate := query.From; !currentDate.After(query.To); currentDate = currentDate.AddDate(0, 0, 1) {
		dates = append(dates, currentDate)
	}

	filter := journals.NewFilter().
		WithAccountIds(query.AccountID).
		WithDates(dates)

	foundJournals, err := handler.journals.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding journals: %w", err)
	}

	filteredJournals := foundJournals[:0]
	for _, journal := range foundJournals {
		if journal == nil || len(journal.Note()) == 0 {
			continue
		}

		filteredJournals = append(filteredJournals, journal)
	}

	sort.Slice(filteredJournals, func(i, j int) bool {
		leftDate := time.Time(filteredJournals[i].Date())
		rightDate := time.Time(filteredJournals[j].Date())

		if leftDate.Equal(rightDate) {
			return filteredJournals[i].CreatedAt().After(filteredJournals[j].CreatedAt())
		}

		return leftDate.After(rightDate)
	})

	return filteredJournals, nil
}
