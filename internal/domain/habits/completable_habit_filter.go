package habits

import (
	"github.com/google/uuid"
)

// CompletableHabitFilter carries optional completable habit lookup criteria.
type CompletableHabitFilter struct {
	AccountIds          []uuid.UUID
	Archived            bool
	Deleted             bool
	CompletableHabitIds []uuid.UUID
}

// NewCompletableHabitFilter returns an empty completable habit filter.
func NewCompletableHabitFilter() CompletableHabitFilter {
	return CompletableHabitFilter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter CompletableHabitFilter) WithAccountIds(accountIds ...uuid.UUID) CompletableHabitFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithArchived returns a filter with the provided archive status.
func (filter CompletableHabitFilter) WithArchived(archived bool) CompletableHabitFilter {
	filter.Archived = archived

	return filter
}

// WithDeleted returns a filter with the provided deletion status.
func (filter CompletableHabitFilter) WithDeleted(deleted bool) CompletableHabitFilter {
	filter.Deleted = deleted

	return filter
}

// WithIds returns a filter with the provided completable habit identifiers.
func (filter CompletableHabitFilter) WithIds(ids ...uuid.UUID) CompletableHabitFilter {
	filter.CompletableHabitIds = append(filter.CompletableHabitIds, ids...)

	return filter
}
