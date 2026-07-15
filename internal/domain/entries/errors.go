package entries

import (
	"github.com/EugeneNail/lifeline/internal/domain"
)

var ErrDateIsOccupied = domain.NewErrorf("date is occupied by another entry")
