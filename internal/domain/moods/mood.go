package moods

import "github.com/EugeneNail/lifeline/internal/domain"

// Mood is a discrete daily mood value.
type Mood int

const (
	Awful Mood = 1
	Bad   Mood = 2
	Okay  Mood = 3
	Good  Mood = 4
	Great Mood = 5
)

var moodLabels = map[Mood]string{
	Awful: "Awful",
	Bad:   "Bad",
	Okay:  "Okay",
	Good:  "Good",
	Great: "Great",
}

// IsValid reports whether the mood value is one of the supported domain values.
func (mood Mood) IsValid() bool {
	_, ok := moodLabels[mood]
	return ok
}

// New returns a validated mood value or an error when rawMood is outside the allowed range.
func New(rawMood int) (Mood, error) {
	if rawMood < int(Awful) || rawMood > int(Great) {
		return 0, domain.NewViolationf("mood must be in range between %d and %d", Awful, Great)
	}

	return Mood(rawMood), nil
}
