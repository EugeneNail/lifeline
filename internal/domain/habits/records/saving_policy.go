package records

import (
	"github.com/EugeneNail/lifeline/internal/domain"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// SavingPolicy validates whether a record can be saved for the given account and date.
type SavingPolicy struct{}

// NewSavingPolicy returns a saving policy instance.
func NewSavingPolicy() *SavingPolicy {
	return &SavingPolicy{}
}

// Check returns nil when the target belongs to the account and is not deleted or archived after the given date, or a domain error otherwise.
func (policy *SavingPolicy) Check(accountID uuid.UUID, date Date, target SavingTarget) domain.Error {
	if target.AccountId() != accountID {
		return habits.ErrHabitBelongsToAnotherUser
	}

	recordDate := time.Time(date)

	if deletedAt := target.DeletedAt(); deletedAt != nil && deletedAt.After(recordDate) {
		return habits.ErrHabitIsDeleted
	}

	if archivedAt := target.ArchivedAt(); archivedAt != nil && archivedAt.After(recordDate) {
		return habits.ErrHabitIsArchived
	}

	return nil
}
