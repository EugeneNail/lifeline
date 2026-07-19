package mood

import "github.com/EugeneNail/lifeline/internal/domain"

// Mood is a discrete daily mood value.
type Mood int

const (
	MoodAwful Mood = 1
	MoodBad   Mood = 2
	MoodOkay  Mood = 3
	MoodGood  Mood = 4
	MoodGreat Mood = 5
)

var moodLabels = map[Mood]string{
	MoodAwful: "Awful",
	MoodBad:   "Bad",
	MoodOkay:  "Okay",
	MoodGood:  "Good",
	MoodGreat: "Great",
}

// New returns a validated mood value or an error when rawMood is outside the allowed range.
func New(rawMood int) (Mood, error) {
	if rawMood < int(MoodAwful) || rawMood > int(MoodGreat) {
		return 0, domain.NewViolationf("mood must be in range between %d and %d", MoodAwful, MoodGreat)
	}

	return Mood(rawMood), nil
}

// String returns the human-readable label for the mood.
func (mood Mood) String() string {
	return moodLabels[mood]
}
