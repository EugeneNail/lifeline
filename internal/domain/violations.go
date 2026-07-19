package domain

// Violations collects domain validation violations keyed by field name.
type Violations struct {
	violations map[string]Violation
}

// NewViolations returns an empty domain validation error collection.
func NewViolations() Violations {
	return Violations{
		violations: make(map[string]Violation),
	}
}

// Error returns a generic validation failure message.
func (violations Violations) Error() string {
	return "domain validation failed"
}

// Add stores a domain validation error under the provided field.
func (violations Violations) Add(field string, err Violation) {
	violations.violations[field] = err
}

// HasViolations reports whether at least one domain validation error has been collected.
func (violations Violations) HasViolations() bool {
	return len(violations.violations) > 0
}

// Violations returns the collected domain validation messages keyed by field name.
func (violations Violations) Violations() map[string]string {
	vltns := make(map[string]string, len(violations.violations))
	for field, err := range violations.violations {
		vltns[field] = err.Error()
	}

	return vltns
}
