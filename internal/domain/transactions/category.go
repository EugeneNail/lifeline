package transactions

import "github.com/EugeneNail/lifeline/internal/domain"

// Category represents a transaction category.
type Category int

const (
	Bills Category = 1 + iota
	Food
	Transport
	Household
	Entertainment
	PersonalItems
	Health
	Work
	Debt
	Investments
	Gifts
	Other
)

var categoryNames = map[Category]string{
	Bills:         "Bills",
	Food:          "Food",
	Transport:     "Transport",
	Household:     "Household",
	Entertainment: "Entertainment",
	PersonalItems: "Personal items",
	Health:        "Health",
	Work:          "Work",
	Debt:          "Debt",
	Investments:   "Investments",
	Gifts:         "Gifts",
	Other:         "Other",
}

// IsValid reports whether the category is one of the supported domain values.
func (category Category) IsValid() bool {
	_, ok := categoryNames[category]

	return ok
}

// NewCategory returns a validated category or a violation when the raw value is outside the supported range.
func NewCategory(rawCategory int) (Category, domain.Violation) {
	category := Category(rawCategory)
	if !category.IsValid() {
		return 0, domain.NewViolationf("category must be in range between %d and %d", Bills, Other)
	}

	return category, nil
}

// String returns the human-readable category name.
func (category Category) String() string {
	name, ok := categoryNames[category]
	if !ok {
		return "Unknown"
	}

	return name
}
