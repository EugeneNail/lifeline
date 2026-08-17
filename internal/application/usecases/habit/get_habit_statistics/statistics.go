package get_habit_statistics

import (
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// Node represents a habit value on a specific date.
type Node struct {
	Date  records.Date
	Value float32
}

// Series represents daily values and range metadata for a habit.
type Series struct {
	HabitID  uuid.UUID
	Nodes    []Node
	MinValue float32
	MaxValue float32
}

// Result carries the habit statistics returned by the use case.
type Result struct {
	MeasurableSeries  []Series
	TimeSeries        []Series
	CompletableSeries []Series
}
