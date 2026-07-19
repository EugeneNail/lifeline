package domain

import (
	"time"
)

// Date is a journal day truncated to calendar precision.
type Date time.Time

// NewDate returns a validated journal date or a violation when the date is invalid.
func NewDate(raw time.Time) (Date, error) {
	if raw.IsZero() {
		return Date{}, NewViolation("date is empty")
	}

	date := raw.Truncate(time.Hour * 24)
	minDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, date.Location())
	maxDate := time.Date(2099, time.January, 1, 0, 0, 0, 0, date.Location())

	if date.Before(minDate) || date.After(maxDate) {
		return Date{}, NewViolation("date must be between 2000-01-01 and 2099-01-01")
	}

	return Date(date), nil
}
