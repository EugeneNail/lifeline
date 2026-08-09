package journals

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// Filter carries optional journal lookup criteria.
type Filter struct {
	AccountIds []auth.ID
	Dates      []time.Time
}

// NewFilter returns an empty journal filter.
func NewFilter() Filter {
	return Filter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter Filter) WithAccountIds(accountIds ...auth.ID) Filter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to day precision.
func (filter Filter) WithDates(dates []time.Time) Filter {
	filter.Dates = append(filter.Dates, truncateJournalDates(dates)...)

	return filter
}

// truncateJournalDates returns the provided dates normalized to YYYY-MM-DD.
func truncateJournalDates(dates []time.Time) []time.Time {
	truncated := make([]time.Time, 0, len(dates))
	for _, date := range dates {
		truncated = append(truncated, time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()))
	}

	return truncated
}
