package habits

import (
	"time"

	"github.com/google/uuid"
)

// ModificationTarget exposes the common state required to validate habit modifications.
type ModificationTarget interface {
	ArchivedAt() *time.Time
	DeletedAt() *time.Time
	AccountId() uuid.UUID
}
