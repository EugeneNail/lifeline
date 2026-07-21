package domain

// Violations collects domain validation violations keyed by field name.
type Violations interface {
	error
	Add(field string, err Violation)
	HasViolations() bool
	Violations() map[string]string
}

type violations struct {
	violations map[string]Violation
}

// NewViolations returns an empty domain validation error collection.
func NewViolations() Violations {
	return violations{
		violations: make(map[string]Violation),
	}
}

// Error returns a generic validation failure message.
func (violations violations) Error() string {
	return "domain validation failed"
}

// Add stores a domain validation error under the provided field.
func (violations violations) Add(field string, err Violation) {
	violations.violations[field] = err
}

// HasViolations reports whether at least one domain validation error has been collected.
func (violations violations) HasViolations() bool {
	return len(violations.violations) > 0
}

// Violations returns the collected domain validation messages keyed by field name.
func (violations violations) Violations() map[string]string {
	vltns := make(map[string]string, len(violations.violations))
	for field, err := range violations.violations {
		vltns[field] = err.Error()
	}

	return vltns
}
