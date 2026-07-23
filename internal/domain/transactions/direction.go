package transactions

import "github.com/EugeneNail/lifeline/internal/domain"

type Direction int

const (
	Expense = 1
	Income  = 2
)

var values = map[Direction]Direction{
	Expense: Expense,
	Income:  Income,
}

func (direction Direction) IsValid() bool {
	_, ok := values[direction]
	return ok
}

func NewDirection(value int) (Direction, domain.Violation) {
	// TODO rewrite checks of the other enums to make them use minmax variables
	minimum := Expense
	maximum := Income
	if value < minimum || value > maximum {
		return 0, domain.NewViolationf("direction must be in range between %d and %d", minimum, maximum)
	}

	return Direction(value), nil
}
