package journal

import (
	"github.com/EugeneNail/lifeline/internal/domain"
)

var ErrDateIsOccupied = domain.NewViolationf("date is occupied by another journal")
