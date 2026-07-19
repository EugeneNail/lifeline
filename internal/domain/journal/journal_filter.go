package journal

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/google/uuid"
)

// JournalFilter carries optional journal lookup criteria.
type JournalFilter struct {
	AccountIds []auth.ID
	Dates      []time.Time
	Ids        []uuid.UUID
}

// NewJournalFilter returns an empty journal filter.
func NewJournalFilter() JournalFilter {
	return JournalFilter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter JournalFilter) WithAccountIds(accountIds ...auth.ID) JournalFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to day precision.
func (filter JournalFilter) WithDates(dates ...time.Time) JournalFilter {
	for _, date := range dates {
		filter.Dates = append(filter.Dates, date.Truncate(time.Hour*24))
	}

	return filter
}

// WithIds returns a filter with the provided journal identifiers.
func (filter JournalFilter) WithIds(ids ...uuid.UUID) JournalFilter {
	filter.Ids = append(filter.Ids, ids...)

	return filter
}
