package journal

import "github.com/EugeneNail/lifeline/internal/domain"

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

func NewMood(rawMood int) (Mood, error) {
	if rawMood < int(MoodAwful) || rawMood > int(MoodGreat) {
		return 0, domain.NewViolationf("mood must be in range between %d and %d", MoodAwful, MoodGreat)
	}

	return Mood(rawMood), nil
}
