package habits

import (
	"github.com/EugeneNail/lifeline/internal/domain"
)

var ErrHabitLimitExceeded = domain.NewViolationf("habit limit is exceeded")

// ErrHabitNotFound reports that a habit with the requested identifier could not be found.
var ErrHabitNotFound = domain.NewViolation("habit not found")

// ErrHabitBelongsToAnotherUser reports that the habit is owned by a different account.
var ErrHabitBelongsToAnotherUser = domain.NewViolation("habit belongs to another user")

// ErrHabitIsArchived reports that an archived habit cannot be modified.
var ErrHabitIsArchived = domain.NewViolation("habit is archived")

// ErrHabitIsDeleted reports that a deleted habit cannot be modified.
var ErrHabitIsDeleted = domain.NewViolation("habit is deleted")
