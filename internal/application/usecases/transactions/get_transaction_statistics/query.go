package get_transaction_statistics

import (
	"time"

	"github.com/google/uuid"
)

// Query carries the account identifier, date range, and optional category filters required to load transaction statistics.
type Query struct {
	AccountID  uuid.UUID
	From       time.Time
	To         time.Time
	Categories []int
}

// NewQuery returns a query with the provided account identifier, normalized date range, and transaction categories.
func NewQuery(accountID uuid.UUID, from time.Time, to time.Time, categories ...int) Query {
	return Query{
		AccountID:  accountID,
		From:       toDate(from),
		To:         toDate(to),
		Categories: append([]int(nil), categories...),
	}
}

// toDate returns the provided time normalized to YYYY-MM-DD.
func toDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
