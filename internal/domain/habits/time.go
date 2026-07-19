package habits

import (
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
)

const (
	timeMinMinute = 0
	timeMaxMinute = 1439
)

// Time represents a time of day as minutes from 00:00 without date information.
type Time int

// NewTime returns a time of day or a domain error when the minute is outside the supported range.
func NewTime(value int) (Time, error) {
	if value < timeMinMinute || value > timeMaxMinute {
		return 0, domain.NewViolationf("time must be between %d and %d minutes", timeMinMinute, timeMaxMinute)
	}

	return Time(value), nil
}

// Hours returns the hour component of the time of day.
func (value Time) Hours() int {
	return int(value) / 60
}

// Minutes returns the minute component of the time of day.
func (value Time) Minutes() int {
	return int(value) % 60
}

// Value returns the time of day as minutes from 00:00.
func (value Time) Value() int {
	return int(value)
}

// String returns the time of day formatted as HH:MM.
func (value Time) String() string {
	return fmt.Sprintf("%02d:%02d", value.Hours(), value.Minutes())
}
