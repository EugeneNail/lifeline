package records

import "time"

// TODO remove and replace with domain.Date
// Date represents a calendar day truncated to 00:00.
type Date time.Time

// NewDate returns a date truncated to a whole day.
func NewDate(rawDate time.Time) Date {
	return Date(rawDate.Truncate(time.Hour * 24))
}
