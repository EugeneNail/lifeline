package transactions

import (
	"github.com/EugeneNail/lifeline/internal/domain"
	"math"
)

const (
	MoneyMin = 0.01
	MoneyMax = 1_000_000_000.0
)

type Money float32

func NewMoney(value float32) (Money, domain.Violation) {
	if value < 0 {
		value = value * -1
	}

	if value < MoneyMin || value > MoneyMax {
		return 0, domain.NewViolationf("value must be between %.2f and %f", MoneyMin, MoneyMax)
	}

	// Rounds to 2 decimal places.
	// 1.625 => 162.5 => 163 => 1.63
	rounded := math.Round(float64(value)*100) / 100

	return Money(float32(rounded)), nil
}
