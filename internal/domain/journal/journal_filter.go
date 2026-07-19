package journal

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/google/uuid"
)

// Filter carries optional journal lookup criteria.
type Filter struct {
	AccountIds []auth.ID
	Dates      []time.Time
	Ids        []uuid.UUID
}

// NewFilter returns an empty journal filter.
func NewFilter() Filter {
	return Filter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter Filter) WithAccountIds(accountIds ...auth.ID) Filter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to day precision.
func (filter Filter) WithDates(dates ...time.Time) Filter {
	for _, date := range dates {
		filter.Dates = append(filter.Dates, date.Truncate(time.Hour*24))
	}

	return filter
}

// WithIds returns a filter with the provided journal identifiers.
func (filter Filter) WithIds(ids ...uuid.UUID) Filter {
	filter.Ids = append(filter.Ids, ids...)

	return filter
}
