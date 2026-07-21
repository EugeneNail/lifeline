package transactions

import (
	"time"

	"github.com/google/uuid"
)

// TransactionFilter carries optional transaction lookup criteria.
type TransactionFilter struct {
	TransactionIds []uuid.UUID
	AccountIds     []uuid.UUID
	Dates          []time.Time
}

// NewTransactionFilter returns an empty transaction filter.
func NewTransactionFilter() TransactionFilter {
	return TransactionFilter{}
}

// WithIds returns a filter with the provided transaction identifiers.
func (filter TransactionFilter) WithIds(ids ...uuid.UUID) TransactionFilter {
	filter.TransactionIds = append(filter.TransactionIds, ids...)

	return filter
}

// WithAccountIds returns a filter with the provided account identifiers.
func (filter TransactionFilter) WithAccountIds(accountIds ...uuid.UUID) TransactionFilter {
	filter.AccountIds = append(filter.AccountIds, accountIds...)

	return filter
}

// WithDates returns a filter with the provided dates truncated to whole days.
func (filter TransactionFilter) WithDates(dates ...time.Time) TransactionFilter {
	for _, date := range dates {
		filter.Dates = append(filter.Dates, date.Truncate(time.Hour*24))
	}

	return filter
}
