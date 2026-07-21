package transactions

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
)

type Transaction struct {
	id          uuid.UUID
	money       Money
	date        domain.Date
	category    Category
	description Description
	accountId   uuid.UUID
	createdAt   time.Time
	updatedAt   time.Time
}

// New returns a transaction with trusted values.
func New(money Money, date domain.Date, category Category, description Description, accountId uuid.UUID) *Transaction {
	now := time.Now()

	return &Transaction{
		id:          generateId(),
		money:       money,
		date:        date,
		category:    category,
		description: description,
		accountId:   accountId,
		createdAt:   now,
		updatedAt:   now,
	}
}

// NewFromRaw returns a transaction with validated fields or domain validation violations.
func NewFromRaw(rawMoney float32, rawDate time.Time, rawCategory int, rawDescription string, accountId uuid.UUID) (*Transaction, domain.Violations) {
	violations := domain.NewViolations()

	money, violation := NewMoney(rawMoney)
	if violation != nil {
		violations.Add("money", violation)
	}

	date, violation := domain.NewDate(rawDate)
	if violation != nil {
		violations.Add("date", violation)
	}

	category, violation := NewCategory(rawCategory)
	if violation != nil {
		violations.Add("category", violation)
	}

	description, violation := NewDescription(rawDescription)
	if violation != nil {
		violations.Add("description", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	return New(money, date, category, description, accountId), nil
}

// Restore returns a transaction reconstructed from persisted primitive values without validating or changing them.
func Restore(id uuid.UUID, money float32, date time.Time, category int, description string, accountId uuid.UUID, createdAt time.Time, updatedAt time.Time) *Transaction {
	return &Transaction{
		id:          id,
		money:       Money(money),
		date:        domain.Date(date),
		category:    Category(category),
		description: Description(description),
		accountId:   accountId,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// generateId returns a new UUIDv7 identifier for a transaction.
func generateId() uuid.UUID {
	return uuid.UUID(uuidv7.New())
}

// ID returns the transaction identifier.
func (transaction *Transaction) ID() uuid.UUID {
	return transaction.id
}

// Id returns the transaction identifier.
func (transaction *Transaction) Id() uuid.UUID {
	return transaction.ID()
}

// Money returns the transaction amount.
func (transaction *Transaction) Money() Money {
	return transaction.money
}

// ChangeMoney updates the transaction amount and refreshes the modification timestamp.
func (transaction *Transaction) ChangeMoney(money Money) {
	transaction.money = money
	transaction.updatedAt = time.Now()
}

// Date returns the transaction date.
func (transaction *Transaction) Date() domain.Date {
	return transaction.date
}

// ChangeDate updates the transaction date and refreshes the modification timestamp.
func (transaction *Transaction) ChangeDate(date domain.Date) {
	transaction.date = date
	transaction.updatedAt = time.Now()
}

// Category returns the transaction category.
func (transaction *Transaction) Category() Category {
	return transaction.category
}

// ChangeCategory updates the transaction category and refreshes the modification timestamp.
func (transaction *Transaction) ChangeCategory(category Category) {
	transaction.category = category
	transaction.updatedAt = time.Now()
}

// Description returns the transaction description.
func (transaction *Transaction) Description() Description {
	return transaction.description
}

// ChangeDescription updates the transaction description and refreshes the modification timestamp.
func (transaction *Transaction) ChangeDescription(description Description) {
	transaction.description = description
	transaction.updatedAt = time.Now()
}

// AccountId returns the identifier of the account that owns the transaction.
func (transaction *Transaction) AccountId() uuid.UUID {
	return transaction.accountId
}

// CreatedAt returns the timestamp when the transaction was created.
func (transaction *Transaction) CreatedAt() time.Time {
	return transaction.createdAt
}

// UpdatedAt returns the timestamp when the transaction was last updated.
func (transaction *Transaction) UpdatedAt() time.Time {
	return transaction.updatedAt
}
