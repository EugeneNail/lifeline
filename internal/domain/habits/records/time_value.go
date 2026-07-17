package records

import (
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
)

const (
	timeValueMinMinute = 0
	timeValueMaxMinute = 1439
)

// TimeValue represents a time of day as minutes from 00:00 without date information.
type TimeValue int

// NewTimeValue returns a time value or a domain error when the minute is outside the supported range.
func NewTimeValue(rawValue int) (TimeValue, error) {
	if rawValue < timeValueMinMinute || rawValue > timeValueMaxMinute {
		return 0, domain.NewErrorf("time value must be between %d and %d minutes", timeValueMinMinute, timeValueMaxMinute)
	}

	return TimeValue(rawValue), nil
}

// Hours returns the hour component of the time value.
func (value TimeValue) Hours() int {
	return int(value) / 60
}

// Minutes returns the minute component of the time value.
func (value TimeValue) Minutes() int {
	return int(value) % 60
}

// Value returns the time value as minutes from 00:00.
func (value TimeValue) Value() int {
	return int(value)
}

// String returns the time value formatted as HH:MM.
func (value TimeValue) String() string {
	return fmt.Sprintf("%02d:%02d", value.Hours(), value.Minutes())
}
