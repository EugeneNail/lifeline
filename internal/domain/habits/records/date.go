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

func truncateRecordDates(dates ...time.Time) []time.Time {
	truncated := make([]time.Time, 0, len(dates))
	for _, date := range dates {
		truncated = append(truncated, date.Truncate(time.Hour*24))
	}

	return truncated
}
