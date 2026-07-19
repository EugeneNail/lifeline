package habits

import (
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
)

// ModificationPolicy validates whether a habit can be modified by the given account.
type ModificationPolicy struct{}

// NewModificationPolicy returns a modification policy instance.
func NewModificationPolicy() *ModificationPolicy {
	return &ModificationPolicy{}
}

// Check returns nil when the target belongs to the account and is neither archived nor deleted, or a domain error otherwise.
func (policy *ModificationPolicy) Check(accountId uuid.UUID, target ModificationTarget) domain.Violation {
	if target.AccountId() != accountId {
		return ErrHabitBelongsToAnotherUser
	}

	if target.ArchivedAt() != nil {
		return ErrHabitIsArchived
	}

	if target.DeletedAt() != nil {
		return ErrHabitIsDeleted
	}

	return nil
}
