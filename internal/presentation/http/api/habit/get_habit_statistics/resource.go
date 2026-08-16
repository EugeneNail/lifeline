package get_habit_statistics

// Output represents the JSON response body for habit statistics.
type Output struct {
	MeasurableSeries  []MeasurableHabitSeries `json:"measurableSeries"`
	TimeSeries        []TimeHabitSeries       `json:"timeSeries"`
	CompletableSeries []CompletableSeries     `json:"completableSeries"`
}

// MeasurableHabitSeries represents the public series data for a measurable habit.
type MeasurableHabitSeries struct {
	HabitID  string  `json:"habitId"`
	Nodes    []Node  `json:"nodes"`
	MaxValue float32 `json:"maxValue"`
}

// TimeHabitSeries represents the public series data for a time habit.
type TimeHabitSeries struct {
	HabitID  string  `json:"habitId"`
	Nodes    []Node  `json:"nodes"`
	MaxValue float32 `json:"maxValue"`
}

// CompletableSeries represents the public series data for a completable habit.
type CompletableSeries struct {
	HabitID string `json:"habitId"`
	Nodes   []Node `json:"nodes"`
}

// Node represents a public habit statistics value for a date.
type Node struct {
	Date  string  `json:"date"`
	Value float32 `json:"value"`
}
