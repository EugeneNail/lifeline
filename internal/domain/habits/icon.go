package habits

import "github.com/EugeneNail/lifeline/internal/domain"

// Icon represents a habit icon.
type Icon int

const (
	IconHeart Icon = 1
	IconBook  Icon = 2
	IconRun   Icon = 3
	IconWater Icon = 4
	IconSleep Icon = 5
)

var iconLabels = map[Icon]string{
	IconHeart: "Heart",
	IconBook:  "Book",
	IconRun:   "Run",
	IconWater: "Water",
	IconSleep: "Sleep",
}

// NewIcon returns an icon enum value or an error when the raw value is unsupported.
func NewIcon(rawIcon int) (Icon, error) {
	icon := Icon(rawIcon)
	if !icon.IsValid() {
		return 0, domain.NewErrorf("icon must be in range between %d and %d", IconHeart, IconSleep)
	}

	return icon, nil
}

// IsValid reports whether the icon is one of the supported enum values.
func (icon Icon) IsValid() bool {
	_, ok := iconLabels[icon]

	return ok
}

// String returns the human-readable icon label.
func (icon Icon) String() string {
	label, ok := iconLabels[icon]
	if !ok {
		return "Unknown"
	}

	return label
}
