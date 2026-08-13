package get_habit_statistics

import (
	"time"

	"github.com/google/uuid"
)

// Query carries the account identifier and date range required to load habit statistics.
type Query struct {
	AccountID uuid.UUID
	From      time.Time
	To        time.Time
}

// NewQuery returns a query with the provided account identifier and normalized date range.
func NewQuery(accountID uuid.UUID, from time.Time, to time.Time) Query {
	return Query{
		AccountID: accountID,
		From:      toDate(from),
		To:        toDate(to),
	}
}

// toDate returns the provided time normalized to YYYY-MM-DD.
func toDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
