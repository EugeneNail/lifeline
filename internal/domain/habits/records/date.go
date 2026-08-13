package records

import "time"

// TODO remove and replace with domain.Date
// Date represents a calendar day truncated to 00:00.
type Date time.Time

// NewDate returns a calendar date normalized to midnight UTC.
func NewDate(rawDate time.Time) Date {
	return Date(time.Date(
		rawDate.Year(),
		rawDate.Month(),
		rawDate.Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	))
}
