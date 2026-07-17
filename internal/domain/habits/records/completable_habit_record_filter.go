package records

import (
	"time"

	"github.com/google/uuid"
)

// CompletableHabitRecordFilter carries optional completable habit record lookup criteria.
type CompletableHabitRecordFilter struct {
	CompletableHabitRecordIds []uuid.UUID
	Dates                     []time.Time
	AccountIds                []uuid.UUID
}

// NewCompletableHabitRecordFilter returns an empty completable habit record filter.
func NewCompletableHabitRecordFilter() CompletableHabitRecordFilter {
	return CompletableHabitRecordFilter{}
}

// WithIds returns a filter with the provided completable habit record identifiers.
func (filter CompletableHabitRecordFilter) WithIds(ids ...uuid.UUID) CompletableHabitRecordFilter {
	filter.CompletableHabitRecordIds = append(filter.CompletableHabitRecordIds, ids...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to whole days.
func (filter CompletableHabitRecordFilter) WithDates(dates ...time.Time) CompletableHabitRecordFilter {
	filter.Dates = append(filter.Dates, truncateRecordDates(dates...)...)

	return filter
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter CompletableHabitRecordFilter) WithAccountIds(accountIds ...uuid.UUID) CompletableHabitRecordFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}
