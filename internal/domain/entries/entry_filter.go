package entries

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// EntryFilter carries optional entry lookup criteria.
type EntryFilter struct {
	AccountIds []auth.ID
	Dates      []time.Time
	Ids        []ID
}

// NewEntryFilter returns an empty entry filter.
func NewEntryFilter() EntryFilter {
	return EntryFilter{}
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter EntryFilter) WithAccountIds(accountIds ...auth.ID) EntryFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to day precision.
func (filter EntryFilter) WithDates(dates ...time.Time) EntryFilter {
	for _, date := range dates {
		filter.Dates = append(filter.Dates, date.Truncate(time.Hour*24))
	}

	return filter
}

// WithIds returns a filter with the provided entry identifiers.
func (filter EntryFilter) WithIds(ids ...ID) EntryFilter {
	filter.Ids = append(filter.Ids, ids...)

	return filter
}
