package get_habit_statistics

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// Usecase executes the get-habit-statistics scenario.
type Usecase struct {
	completableHabits       habits.CompletableHabitRepository
	completableHabitRecords records.CompletableHabitRecordRepository
	measurableHabits        habits.MeasurableHabitRepository
	measurableHabitRecords  records.MeasurableHabitRecordRepository
	timeHabits              habits.TimeHabitRepository
	timeHabitRecords        records.TimeHabitRecordRepository
}

// NewUsecase returns a get-habit-statistics use case configured with the habit and record repositories or an error when a dependency is missing.
func NewUsecase(
	measurableHabits habits.MeasurableHabitRepository,
	measurableHabitRecords records.MeasurableHabitRecordRepository,
	completableHabits habits.CompletableHabitRepository,
	completableHabitRecords records.CompletableHabitRecordRepository,
	timeHabits habits.TimeHabitRepository,
	timeHabitRecords records.TimeHabitRecordRepository,
) (*Usecase, error) {
	if measurableHabits == nil {
		return nil, fmt.Errorf("get_habit_statistics usecase requires a measurable habit repository")
	}

	if measurableHabitRecords == nil {
		return nil, fmt.Errorf("get_habit_statistics usecase requires a measurable habit record repository")
	}

	if completableHabits == nil {
		return nil, fmt.Errorf("get_habit_statistics usecase requires a completable habit repository")
	}

	if completableHabitRecords == nil {
		return nil, fmt.Errorf("get_habit_statistics usecase requires a completable habit record repository")
	}

	if timeHabits == nil {
		return nil, fmt.Errorf("get_habit_statistics usecase requires a time habit repository")
	}

	if timeHabitRecords == nil {
		return nil, fmt.Errorf("get_habit_statistics usecase requires a time habit record repository")
	}

	return &Usecase{
		completableHabits:       completableHabits,
		completableHabitRecords: completableHabitRecords,
		measurableHabits:        measurableHabits,
		measurableHabitRecords:  measurableHabitRecords,
		timeHabits:              timeHabits,
		timeHabitRecords:        timeHabitRecords,
	}, nil
}

// Handle returns measurable, time, and completable habit heatmaps for the requested range or an error when the range is invalid or data cannot be loaded.
func (usecase *Usecase) Handle(ctx context.Context, query Query) (Result, error) {
	if query.From.After(query.To) {
		return Result{}, ErrInvalidDateRange
	}

	foundMeasurableHabits, err := usecase.measurableHabits.FindMany(
		ctx,
		habits.NewMeasurableHabitFilter().
			WithAccountIds(query.AccountID).
			WithArchived(false).
			WithDeleted(false),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding active measurable habits for account id %q: %w", query.AccountID, err)
	}

	dates := buildDateRange(query.From, query.To)
	foundMeasurableRecords, err := usecase.measurableHabitRecords.FindMany(
		ctx,
		records.NewMeasurableHabitRecordFilter().
			WithAccountIds(query.AccountID).
			WithDates(dates...),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding measurable habit records for account id %q between %q and %q: %w", query.AccountID, query.From, query.To, err)
	}

	foundCompletableHabits, err := usecase.completableHabits.FindMany(
		ctx,
		habits.NewCompletableHabitFilter().
			WithAccountIds(query.AccountID).
			WithArchived(false).
			WithDeleted(false),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding active completable habits for account id %q: %w", query.AccountID, err)
	}

	foundCompletableRecords, err := usecase.completableHabitRecords.FindMany(
		ctx,
		records.NewCompletableHabitRecordFilter().
			WithAccountIds(query.AccountID).
			WithDates(dates...),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding completable habit records for account id %q between %q and %q: %w", query.AccountID, query.From, query.To, err)
	}

	foundTimeHabits, err := usecase.timeHabits.FindMany(
		ctx,
		habits.NewTimeHabitFilter().
			WithAccountIds(query.AccountID).
			WithArchived(false).
			WithDeleted(false),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding active time habits for account id %q: %w", query.AccountID, err)
	}

	foundTimeRecords, err := usecase.timeHabitRecords.FindMany(
		ctx,
		records.NewTimeHabitRecordFilter().
			WithAccountIds(query.AccountID).
			WithDates(dates...),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding time habit records for account id %q between %q and %q: %w", query.AccountID, query.From, query.To, err)
	}

	return Result{
		MeasurableHeatmap:  buildMeasurableHeatmaps(foundMeasurableHabits, foundMeasurableRecords, dates),
		TimeHeatmap:        buildTimeHeatmaps(foundTimeHabits, foundTimeRecords, dates),
		CompletableHeatmap: buildCompletableHeatmaps(foundCompletableHabits, foundCompletableRecords, dates),
	}, nil
}

// buildDateRange returns every date in the inclusive range sorted in ascending order.
func buildDateRange(from time.Time, to time.Time) []time.Time {
	dates := make([]time.Time, 0)
	for date := from; !date.After(to); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date)
	}

	sort.Slice(dates, func(i int, j int) bool {
		return dates[i].Before(dates[j])
	})

	return dates
}

// buildMeasurableHeatmaps creates a complete heatmap for every active measurable habit and fills dates without records with zero.
func buildMeasurableHeatmaps(
	foundHabits []*habits.MeasurableHabit,
	foundRecords []*records.MeasurableHabitRecord,
	dates []time.Time,
) []MeasurableHabitHeatmap {
	valuesByHabit := make(map[uuid.UUID]map[records.Date]float32)
	maxValues := make(map[uuid.UUID]float32)
	habitIDs := make([]uuid.UUID, 0, len(foundHabits))

	for _, habit := range foundHabits {
		if habit == nil {
			continue
		}

		habitID := habit.ID()
		valuesByHabit[habitID] = make(map[records.Date]float32)
		habitIDs = append(habitIDs, habitID)
	}

	for _, record := range foundRecords {
		if record == nil {
			continue
		}

		habitID := record.MeasurableHabitId()
		if _, exists := valuesByHabit[habitID]; !exists {
			continue
		}

		value := record.Value().Value()
		valuesByHabit[habitID][record.Date()] = value
		if value > maxValues[habitID] {
			maxValues[habitID] = value
		}
	}

	sort.Slice(habitIDs, func(i int, j int) bool {
		return habitIDs[i].String() < habitIDs[j].String()
	})

	heatmaps := make([]MeasurableHabitHeatmap, 0, len(habitIDs))
	for _, habitID := range habitIDs {
		nodes := make([]HeatmapNode, 0, len(dates))
		for _, date := range dates {
			recordDate := records.NewDate(date)
			nodes = append(nodes, HeatmapNode{
				Date:  recordDate,
				Value: valuesByHabit[habitID][recordDate],
			})
		}

		heatmaps = append(heatmaps, MeasurableHabitHeatmap{
			HabitID:  habitID,
			Nodes:    nodes,
			MaxValue: maxValues[habitID],
		})
	}

	return heatmaps
}

// buildTimeHeatmaps creates a complete heatmap for every active time habit, fills dates without records with zero, and calculates each maximum value.
func buildTimeHeatmaps(
	foundHabits []*habits.TimeHabit,
	foundRecords []*records.TimeHabitRecord,
	dates []time.Time,
) []TimeHabitHeatmap {
	valuesByHabit := make(map[uuid.UUID]map[records.Date]float32)
	maxValues := make(map[uuid.UUID]float32)
	habitIDs := make([]uuid.UUID, 0, len(foundHabits))

	for _, habit := range foundHabits {
		if habit == nil {
			continue
		}

		habitID := habit.ID()
		valuesByHabit[habitID] = make(map[records.Date]float32)
		habitIDs = append(habitIDs, habitID)
	}

	for _, record := range foundRecords {
		if record == nil {
			continue
		}

		habitID := record.TimeHabitId()
		if _, exists := valuesByHabit[habitID]; !exists {
			continue
		}

		value := float32(record.Value().Value())
		valuesByHabit[habitID][record.Date()] = value
		if value > maxValues[habitID] {
			maxValues[habitID] = value
		}
	}

	sort.Slice(habitIDs, func(i int, j int) bool {
		return habitIDs[i].String() < habitIDs[j].String()
	})

	heatmaps := make([]TimeHabitHeatmap, 0, len(habitIDs))
	for _, habitID := range habitIDs {
		nodes := make([]HeatmapNode, 0, len(dates))
		for _, date := range dates {
			recordDate := records.NewDate(date)
			nodes = append(nodes, HeatmapNode{
				Date:  recordDate,
				Value: valuesByHabit[habitID][recordDate],
			})
		}

		heatmaps = append(heatmaps, TimeHabitHeatmap{
			HabitID:  habitID,
			Nodes:    nodes,
			MaxValue: maxValues[habitID],
		})
	}

	return heatmaps
}

// buildCompletableHeatmaps creates a complete heatmap for every active completable habit using zero for no record, one for incomplete, and two for complete.
func buildCompletableHeatmaps(
	foundHabits []*habits.CompletableHabit,
	foundRecords []*records.CompletableHabitRecord,
	dates []time.Time,
) []CompletableHeatmap {
	valuesByHabit := make(map[uuid.UUID]map[records.Date]float32)
	habitIDs := make([]uuid.UUID, 0, len(foundHabits))

	for _, habit := range foundHabits {
		if habit == nil {
			continue
		}

		habitID := habit.ID()
		valuesByHabit[habitID] = make(map[records.Date]float32)
		habitIDs = append(habitIDs, habitID)
	}

	for _, record := range foundRecords {
		if record == nil {
			continue
		}

		habitID := record.CompletableHabitId()
		if _, exists := valuesByHabit[habitID]; !exists {
			continue
		}

		value := float32(1)
		if record.Value() {
			value = 2
		}
		valuesByHabit[habitID][record.Date()] = value
	}

	sort.Slice(habitIDs, func(i int, j int) bool {
		return habitIDs[i].String() < habitIDs[j].String()
	})

	heatmaps := make([]CompletableHeatmap, 0, len(habitIDs))
	for _, habitID := range habitIDs {
		nodes := make([]HeatmapNode, 0, len(dates))
		for _, date := range dates {
			recordDate := records.NewDate(date)
			nodes = append(nodes, HeatmapNode{
				Date:  recordDate,
				Value: valuesByHabit[habitID][recordDate],
			})
		}

		heatmaps = append(heatmaps, CompletableHeatmap{
			HabitID: habitID,
			Nodes:   nodes,
		})
	}

	return heatmaps
}
