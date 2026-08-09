package list_journals

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// Query carries the data required to list journals for an account and date range.
type Query struct {
	AccountID auth.ID
	From      time.Time
	To        time.Time
}

// NewQuery returns a query with the provided account and date range normalized to whole days.
func NewQuery(accountID auth.ID, from time.Time, to time.Time) Query {
	return Query{
		AccountID: accountID,
		From:      toDate(from),
		To:        toDate(to),
	}
}
