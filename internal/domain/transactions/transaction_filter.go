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
	From           *time.Time
	To             *time.Time
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

// WithFrom returns a filter that includes transactions from the provided date inclusive.
func (filter TransactionFilter) WithFrom(from time.Time) TransactionFilter {
	date := from.Truncate(time.Hour * 24)
	filter.From = &date

	return filter
}

// WithTo returns a filter that includes transactions up to the provided date inclusive.
func (filter TransactionFilter) WithTo(to time.Time) TransactionFilter {
	date := to.Truncate(time.Hour * 24)
	filter.To = &date

	return filter
}
