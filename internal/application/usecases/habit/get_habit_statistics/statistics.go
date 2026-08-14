package get_habit_statistics

import (
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// HeatmapNode represents a measurable habit value on a specific date.
type HeatmapNode struct {
	Date  records.Date
	Value float32
}

// MeasurableHabitHeatmap represents daily values and the maximum value for a measurable habit.
type MeasurableHabitHeatmap struct {
	HabitID  uuid.UUID
	Nodes    []HeatmapNode
	MaxValue float32
}

// CompletableHeatmap represents daily completion states for a completable habit.
type CompletableHeatmap struct {
	HabitID uuid.UUID
	Nodes   []HeatmapNode
}

// Result carries the habit statistics returned by the use case.
type Result struct {
	MeasurableHeatmap  []MeasurableHabitHeatmap
	CompletableHeatmap []CompletableHeatmap
}
