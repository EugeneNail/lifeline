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

// MeasurableHabitSeries represents daily values and the maximum value for a measurable habit.
type MeasurableHabitSeries struct {
	HabitID  uuid.UUID
	Nodes    []Node
	MaxValue float32
}

// TimeHabitSeries represents daily time values and the maximum value for a time habit.
type TimeHabitSeries struct {
	HabitID  uuid.UUID
	Nodes    []Node
	MaxValue float32
}

// CompletableSeries represents daily completion states for a completable habit.
type CompletableSeries struct {
	HabitID uuid.UUID
	Nodes   []Node
}

// Result carries the habit statistics returned by the use case.
type Result struct {
	MeasurableSeries  []MeasurableHabitSeries
	TimeSeries        []TimeHabitSeries
	CompletableSeries []CompletableSeries
}
