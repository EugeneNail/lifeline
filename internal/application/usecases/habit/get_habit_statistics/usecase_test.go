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
