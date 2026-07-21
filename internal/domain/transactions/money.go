package transactions

import "github.com/EugeneNail/lifeline/internal/domain"

const (
	MoneyMin = 0.01
	MoneyMax = 1_000_000_000.0
)

type Money float32

func NewMoney(value float32) (Money, domain.Violation) {
	if value < MoneyMin || value > MoneyMax {
		return 0, domain.NewViolationf("value must be between %.2f and %f", MoneyMin, MoneyMax)
	}

	return Money(value), nil
}
