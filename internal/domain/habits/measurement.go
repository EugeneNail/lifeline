package habits

import (
	"unicode/utf8"

	"github.com/EugeneNail/lifeline/internal/domain"
)

const (
	measurementStepMin      = 0.01
	measurementStepMax      = 100000.0
	measurableUnitMinLength = 1
	measurableUnitMaxLength = 8
)

// MeasurementStep represents the smallest measurable habit increment.
type MeasurementStep float32

// NewMeasurementStep returns a measurement step or a domain error when the value is outside the supported range.
func NewMeasurementStep(rawStep float32) (MeasurementStep, error) {
	if rawStep < measurementStepMin || rawStep > measurementStepMax {
		return 0, domain.NewErrorf("step must be between %.2f and %.0f", measurementStepMin, measurementStepMax)
	}

	return MeasurementStep(rawStep), nil
}

// MeasurableUnit represents the name of a measurable habit unit.
type MeasurableUnit string

// NewMeasurableUnit returns a measurement unit name or a domain error when the value violates domain rules.
func NewMeasurableUnit(rawUnit string) (MeasurableUnit, error) {
	length := utf8.RuneCountInString(rawUnit)
	if length < measurableUnitMinLength || length > measurableUnitMaxLength {
		return "", domain.NewErrorf(
			"unit length must be between %d and %d characters",
			measurableUnitMinLength,
			measurableUnitMaxLength,
		)
	}

	for _, character := range rawUnit {
		if character >= 'a' && character <= 'z' {
			continue
		}

		if character == '.' {
			continue
		}

		return "", domain.NewError("unit contains unsupported characters")
	}

	return MeasurableUnit(rawUnit), nil
}
