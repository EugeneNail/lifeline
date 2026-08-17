package get_habit_statistics

// Output represents the JSON response body for habit statistics.
type Output struct {
	MeasurableSeries  []Series `json:"measurableSeries"`
	TimeSeries        []Series `json:"timeSeries"`
	CompletableSeries []Series `json:"completableSeries"`
}

// Series represents the public series data for a habit.
type Series struct {
	HabitID  string  `json:"habitId"`
	Nodes    []Node  `json:"nodes"`
	MinValue float32 `json:"minValue"`
	MaxValue float32 `json:"maxValue"`
}

// Node represents a public habit statistics value for a date.
type Node struct {
	Date  string  `json:"date"`
	Value float32 `json:"value"`
}
