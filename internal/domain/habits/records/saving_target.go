package records

import (
	"time"

	"github.com/google/uuid"
)

// SavingTarget exposes the habit state required to decide whether a record can be saved.
type SavingTarget interface {
	ArchivedAt() *time.Time
	DeletedAt() *time.Time
	AccountId() uuid.UUID
}
