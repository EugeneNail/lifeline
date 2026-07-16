package habits

import (
	"github.com/EugeneNail/lifeline/internal/domain"
)

var ErrHabitLimitExceeded = domain.NewErrorf("habit limit is exceeded")
