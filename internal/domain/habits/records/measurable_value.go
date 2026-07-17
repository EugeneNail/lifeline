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
func NewMeasurableValue(rawValue float32, step habits.MeasurementStep) (MeasurableValue, domain.Error) {
	if rawValue < measurableValueMin || rawValue > measurableValueMax {
		return 0, domain.NewErrorf("measurable value must be between %.2f and %f", measurableValueMin, measurableValueMax)
	}

	value := float64(rawValue)
	stepValue := float64(step)
	// The record value must be an exact multiple of the habit step.
	// Example: if step is 0.25, then 1.00, 1.25, 1.50, and 1.75 are valid.
	// We use math.Mod to compute the remainder after division by step:
	// - remainder == 0 means the value fits the step exactly
	// - a small remainder can still appear because float32/float64 math is not perfectly precise
	// That is why we allow only a tiny epsilon instead of comparing raw floats with strict equality.
	remainder := math.Mod(value, stepValue)
	if remainder != 0 {
		const epsilon = 1e-6
		if remainder > epsilon && stepValue-remainder > epsilon {
			return 0, domain.NewErrorf("measurable value must be a multiple of step %g", step)
		}
	}

	return MeasurableValue(rawValue), nil
}

// Value returns the measurable value as a float32.
func (value MeasurableValue) Value() float32 {
	return float32(value)
}
