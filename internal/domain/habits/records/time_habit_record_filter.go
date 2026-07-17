package records

import (
	"time"

	"github.com/google/uuid"
)

// TimeHabitRecordFilter carries optional time habit record lookup criteria.
type TimeHabitRecordFilter struct {
	TimeHabitRecordIds []uuid.UUID
	Dates              []time.Time
	AccountIds         []uuid.UUID
}

// NewTimeHabitRecordFilter returns an empty time habit record filter.
func NewTimeHabitRecordFilter() TimeHabitRecordFilter {
	return TimeHabitRecordFilter{}
}

// WithIds returns a filter with the provided time habit record identifiers.
func (filter TimeHabitRecordFilter) WithIds(ids ...uuid.UUID) TimeHabitRecordFilter {
	filter.TimeHabitRecordIds = append(filter.TimeHabitRecordIds, ids...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to whole days.
func (filter TimeHabitRecordFilter) WithDates(dates ...time.Time) TimeHabitRecordFilter {
	filter.Dates = append(filter.Dates, truncateRecordDates(dates...)...)

	return filter
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter TimeHabitRecordFilter) WithAccountIds(accountIds ...uuid.UUID) TimeHabitRecordFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}
