package get_habit_statistics

// Output represents the JSON response body for habit statistics.
type Output struct {
	MeasurableHeatmap []MeasurableHabitHeatmap `json:"measurableHeatmap"`
}

// MeasurableHabitHeatmap represents the public heatmap data for a measurable habit.
type MeasurableHabitHeatmap struct {
	HabitID  string        `json:"habitId"`
	Nodes    []HeatmapNode `json:"nodes"`
	MaxValue float32       `json:"maxValue"`
}

// HeatmapNode represents the public measurable habit value for a date.
type HeatmapNode struct {
	Date  string  `json:"date"`
	Value float32 `json:"value"`
}
