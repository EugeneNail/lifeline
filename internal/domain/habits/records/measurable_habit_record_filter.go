package records

import (
	"time"

	"github.com/google/uuid"
)

// MeasurableHabitRecordFilter carries optional measurable habit record lookup criteria.
type MeasurableHabitRecordFilter struct {
	MeasurableHabitRecordIds []uuid.UUID
	Dates                    []time.Time
	AccountIds               []uuid.UUID
}

// NewMeasurableHabitRecordFilter returns an empty measurable habit record filter.
func NewMeasurableHabitRecordFilter() MeasurableHabitRecordFilter {
	return MeasurableHabitRecordFilter{}
}

// WithIds returns a filter with the provided measurable habit record identifiers.
func (filter MeasurableHabitRecordFilter) WithIds(ids ...uuid.UUID) MeasurableHabitRecordFilter {
	filter.MeasurableHabitRecordIds = append(filter.MeasurableHabitRecordIds, ids...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to whole days.
func (filter MeasurableHabitRecordFilter) WithDates(dates ...time.Time) MeasurableHabitRecordFilter {
	filter.Dates = append(filter.Dates, truncateRecordDates(dates...)...)

	return filter
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter MeasurableHabitRecordFilter) WithAccountIds(accountIds ...uuid.UUID) MeasurableHabitRecordFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}
