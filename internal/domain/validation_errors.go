package domain

// ValidationErrors collects domain validation errors keyed by field name.
type ValidationErrors struct {
	errors map[string]Error
}

// NewValidationErrors returns an empty domain validation error collection.
func NewValidationErrors() ValidationErrors {
	return ValidationErrors{
		errors: make(map[string]Error),
	}
}

// Error returns a generic validation failure message.
func (validationErrors ValidationErrors) Error() string {
	return "domain validation failed"
}

// Add stores a domain validation error under the provided field.
func (validationErrors ValidationErrors) Add(field string, err Error) {
	validationErrors.errors[field] = err
}

// HasErrors reports whether at least one domain validation error has been collected.
func (validationErrors ValidationErrors) HasErrors() bool {
	return len(validationErrors.errors) > 0
}

// Errors returns the collected domain validation messages keyed by field name.
func (validationErrors ValidationErrors) Errors() map[string]string {
	errors := make(map[string]string, len(validationErrors.errors))
	for field, err := range validationErrors.errors {
		errors[field] = err.Error()
	}

	return errors
}
