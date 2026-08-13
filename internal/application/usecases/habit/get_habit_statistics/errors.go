package get_habit_statistics

import "errors"

// ErrInvalidDateRange reports that the from date is after the to date.
var ErrInvalidDateRange = errors.New("'from' date must not be after 'to' date")
