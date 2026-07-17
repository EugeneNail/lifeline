package habits

import (
	"github.com/EugeneNail/lifeline/internal/domain"
)

var ErrHabitLimitExceeded = domain.NewErrorf("habit limit is exceeded")

// ErrHabitNotFound reports that a habit with the requested identifier could not be found.
var ErrHabitNotFound = domain.NewError("habit not found")

// ErrHabitBelongsToAnotherUser reports that the habit is owned by a different account.
var ErrHabitBelongsToAnotherUser = domain.NewError("habit belongs to another user")

// ErrHabitIsArchived reports that an archived habit cannot be modified.
var ErrHabitIsArchived = domain.NewError("habit is archived")

// ErrHabitIsDeleted reports that a deleted habit cannot be modified.
var ErrHabitIsDeleted = domain.NewError("habit is deleted")
