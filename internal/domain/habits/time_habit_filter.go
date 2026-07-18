package habits

import (
	"github.com/google/uuid"
)

// TimeHabitFilter carries optional time habit lookup criteria.
type TimeHabitFilter struct {
	AccountIds   []uuid.UUID
	Archived     *bool
	Deleted      *bool
	TimeHabitIds []uuid.UUID
}

// NewTimeHabitFilter returns an empty time habit filter.
func NewTimeHabitFilter() TimeHabitFilter {
	return TimeHabitFilter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter TimeHabitFilter) WithAccountIds(accountIds ...uuid.UUID) TimeHabitFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithArchived returns a filter with the provided archive status.
func (filter TimeHabitFilter) WithArchived(archived bool) TimeHabitFilter {
	value := archived
	filter.Archived = &value

	return filter
}

// WithDeleted returns a filter with the provided deletion status.
func (filter TimeHabitFilter) WithDeleted(deleted bool) TimeHabitFilter {
	value := deleted
	filter.Deleted = &value

	return filter
}

// WithIds returns a filter with the provided time habit identifiers.
func (filter TimeHabitFilter) WithIds(ids ...uuid.UUID) TimeHabitFilter {
	filter.TimeHabitIds = append(filter.TimeHabitIds, ids...)

	return filter
}
