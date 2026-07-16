package habits

import "github.com/google/uuid"

// MeasurableHabitFilter carries optional measurable habit lookup criteria.
type MeasurableHabitFilter struct {
	AccountIds         []uuid.UUID
	Archived           bool
	Deleted            bool
	MeasurableHabitIds []uuid.UUID
}

// NewMeasurableHabitFilter returns an empty measurable habit filter.
func NewMeasurableHabitFilter() MeasurableHabitFilter {
	return MeasurableHabitFilter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter MeasurableHabitFilter) WithAccountIds(accountIds ...uuid.UUID) MeasurableHabitFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)
	return filter
}

// WithArchived returns a filter with the provided archive status.
func (filter MeasurableHabitFilter) WithArchived(archived bool) MeasurableHabitFilter {
	filter.Archived = archived
	return filter
}

// WithDeleted returns a filter with the provided deletion status.
func (filter MeasurableHabitFilter) WithDeleted(deleted bool) MeasurableHabitFilter {
	filter.Deleted = deleted
	return filter
}

// WithIds returns a filter with the provided measurable habit identifiers.
func (filter MeasurableHabitFilter) WithIds(ids ...uuid.UUID) MeasurableHabitFilter {
	filter.MeasurableHabitIds = append(filter.MeasurableHabitIds, ids...)
	return filter
}
