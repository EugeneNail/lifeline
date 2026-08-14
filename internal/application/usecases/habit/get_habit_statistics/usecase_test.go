package get_habit_statistics

import (
	"testing"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// TestBuildMeasurableHeatmapsMatchesDatesAcrossLocations verifies that PostgreSQL dates and query dates resolve to the same heatmap node.
func TestBuildMeasurableHeatmapsMatchesDatesAcrossLocations(t *testing.T) {
	habitID := uuid.MustParse("00000000-0000-7000-8000-000000000001")
	accountID := uuid.MustParse("00000000-0000-7000-8000-000000000002")
	queryDate := time.Date(2026, time.August, 13, 0, 0, 0, 0, time.UTC)
	postgresDate := time.Date(2026, time.August, 13, 0, 0, 0, 0, time.FixedZone("", 0))

	habit := habits.RestoreMeasurableHabit(
		habitID,
		"Water",
		1,
		1,
		"litres",
		queryDate,
		queryDate,
		nil,
		nil,
		accountID,
	)
	record := records.RestoreMeasurableHabitRecord(
		habitID,
		accountID,
		records.NewDate(postgresDate),
		records.MeasurableValue(2.5),
	)

	heatmaps := buildMeasurableHeatmaps(
		[]*habits.MeasurableHabit{habit},
		[]*records.MeasurableHabitRecord{record},
		[]time.Time{queryDate},
	)

	if len(heatmaps) != 1 {
		t.Fatalf("expected one heatmap, got %d", len(heatmaps))
	}

	if len(heatmaps[0].Nodes) != 1 {
		t.Fatalf("expected one heatmap node, got %d", len(heatmaps[0].Nodes))
	}

	if heatmaps[0].Nodes[0].Value != 2.5 {
		t.Fatalf("expected heatmap node value 2.5, got %v", heatmaps[0].Nodes[0].Value)
	}

	if heatmaps[0].MaxValue != 2.5 {
		t.Fatalf("expected maximum value 2.5, got %v", heatmaps[0].MaxValue)
	}
}

// TestBuildCompletableHeatmapsIncludesEveryHabitAndDate verifies that completable heatmaps distinguish missing, incomplete, and complete records.
func TestBuildCompletableHeatmapsIncludesEveryHabitAndDate(t *testing.T) {
	firstHabitID := uuid.MustParse("00000000-0000-7000-8000-000000000001")
	secondHabitID := uuid.MustParse("00000000-0000-7000-8000-000000000002")
	accountID := uuid.MustParse("00000000-0000-7000-8000-000000000003")
	firstDate := time.Date(2026, time.August, 11, 0, 0, 0, 0, time.UTC)
	secondDate := firstDate.AddDate(0, 0, 1)
	thirdDate := firstDate.AddDate(0, 0, 2)

	firstHabit := habits.RestoreCompletableHabit(
		firstHabitID,
		"Gym",
		1,
		firstDate,
		firstDate,
		nil,
		nil,
		accountID,
	)
	secondHabit := habits.RestoreCompletableHabit(
		secondHabitID,
		"Read",
		2,
		firstDate,
		firstDate,
		nil,
		nil,
		accountID,
	)
	incompleteRecord := records.RestoreCompletableHabitRecord(
		firstHabitID,
		accountID,
		records.NewDate(firstDate),
		false,
	)
	completeRecord := records.RestoreCompletableHabitRecord(
		firstHabitID,
		accountID,
		records.NewDate(thirdDate),
		true,
	)

	heatmaps := buildCompletableHeatmaps(
		[]*habits.CompletableHabit{secondHabit, firstHabit},
		[]*records.CompletableHabitRecord{incompleteRecord, completeRecord},
		[]time.Time{firstDate, secondDate, thirdDate},
	)

	if len(heatmaps) != 2 {
		t.Fatalf("expected two heatmaps, got %d", len(heatmaps))
	}

	if heatmaps[0].HabitID != firstHabitID {
		t.Fatalf("expected first heatmap for habit %q, got %q", firstHabitID, heatmaps[0].HabitID)
	}

	assertHeatmapNodeValues(t, heatmaps[0].Nodes, []float32{1, 0, 2})
	assertHeatmapNodeValues(t, heatmaps[1].Nodes, []float32{0, 0, 0})
}

// assertHeatmapNodeValues verifies that heatmap nodes contain the expected values in date order.
func assertHeatmapNodeValues(t *testing.T, nodes []HeatmapNode, expected []float32) {
	t.Helper()

	if len(nodes) != len(expected) {
		t.Fatalf("expected %d nodes, got %d", len(expected), len(nodes))
	}

	for index, expectedValue := range expected {
		if nodes[index].Value != expectedValue {
			t.Errorf("expected node %d value %v, got %v", index, expectedValue, nodes[index].Value)
		}
	}
}
