package records

import (
	"math"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
)

const (
	measurableValueMin = 0.01
	measurableValueMax = 1000000000.0
)

// MeasurableValue represents a numeric value entered for a measurable habit record.
type MeasurableValue float32

// NewMeasurableValue returns a measurable value or a domain error when the value violates domain rules.
func NewMeasurableValue(rawValue float32, step habits.MeasurementStep) (MeasurableValue, domain.Violation) {
	if rawValue < measurableValueMin || rawValue > measurableValueMax {
		return 0, domain.NewViolationf("measurable value must be between %.2f and %f", measurableValueMin, measurableValueMax)
	}

	value := float64(rawValue)
	stepValue := float64(step)
	// The record value must land on one of the step boundaries.
	//
	// Example: if step is 0.1, then 91.7 is valid even when the raw float comes in
	// as 91.69999695 from JSON, because the source of truth is still "91.7".
	//
	// We therefore do not compare the raw float for strict equality. Instead, we:
	// 1. divide the value by step,
	// 2. round to the nearest whole step count,
	// 3. rebuild the value from that count,
	// 4. compare the difference with a small epsilon.
	//
	// If the value is far from the nearest step boundary, it is rejected.
	const epsilon = 1e-5
	nearestStepCount := math.Round(value / stepValue)
	normalizedValue := nearestStepCount * stepValue
	if math.Abs(value-normalizedValue) > epsilon {
		return 0, domain.NewViolationf("measurable value %0.8f must be a multiple of step %g", value, step)
	}

	return MeasurableValue(float32(normalizedValue)), nil
}

// Value returns the measurable value as a float32.
func (value MeasurableValue) Value() float32 {
	return float32(value)
}
