package moods

import (
	"time"

	"github.com/google/uuid"
)

// RecordFilter carries optional mood record lookup criteria.
type RecordFilter struct {
	AccountIds []uuid.UUID
	Dates      []time.Time
}

// NewRecordFilter returns an empty mood record filter.
func NewRecordFilter() RecordFilter {
	return RecordFilter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter RecordFilter) WithAccountIds(accountIds ...uuid.UUID) RecordFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to whole days.
func (filter RecordFilter) WithDates(dates ...time.Time) RecordFilter {
	filter.Dates = append(filter.Dates, truncateMoodRecordDates(dates...)...)

	return filter
}

func truncateMoodRecordDates(dates ...time.Time) []time.Time {
	truncated := make([]time.Time, 0, len(dates))
	for _, date := range dates {
		truncated = append(truncated, date.Truncate(time.Hour*24))
	}

	return truncated
}
