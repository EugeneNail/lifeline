package get_habit_statistics

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	application "github.com/EugeneNail/lifeline/internal/application/usecases/habit/get_habit_statistics"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the get-habit-statistics use case to the HTTP transport.
type Handler struct {
	usecase  *application.Usecase
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-habit-statistics use case.
func NewHandler(usecase *application.Usecase, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Handle loads habit statistics for the requested date range and returns an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	from, err := time.Parse(time.DateOnly, request.URL.Query().Get("from"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing a 'from' date: %w", err)
	}

	to, err := time.Parse(time.DateOnly, request.URL.Query().Get("to"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing a 'to' date: %w", err)
	}

	result, err := handler.usecase.Handle(request.Context(), application.NewQuery(accountID.Uuid(), from, to))
	if err != nil {
		if errors.Is(err, application.ErrInvalidDateRange) {
			return http.StatusBadRequest, err
		}

		return http.StatusInternalServerError, fmt.Errorf("handling GetHabitStatistics query: %w", err)
	}

	return http.StatusOK, Output{
		MeasurableHeatmap:  mapMeasurableHeatmaps(result.MeasurableHeatmap),
		CompletableHeatmap: mapCompletableHeatmaps(result.CompletableHeatmap),
	}
}

// mapMeasurableHeatmaps converts application heatmaps to transport output.
func mapMeasurableHeatmaps(heatmaps []application.MeasurableHabitHeatmap) []MeasurableHabitHeatmap {
	output := make([]MeasurableHabitHeatmap, 0, len(heatmaps))
	for _, heatmap := range heatmaps {
		output = append(output, MeasurableHabitHeatmap{
			HabitID:  heatmap.HabitID.String(),
			Nodes:    mapHeatmapNodes(heatmap.Nodes),
			MaxValue: heatmap.MaxValue,
		})
	}

	return output
}

// mapCompletableHeatmaps converts application completable heatmaps to transport output.
func mapCompletableHeatmaps(heatmaps []application.CompletableHeatmap) []CompletableHeatmap {
	output := make([]CompletableHeatmap, 0, len(heatmaps))
	for _, heatmap := range heatmaps {
		output = append(output, CompletableHeatmap{
			HabitID: heatmap.HabitID.String(),
			Nodes:   mapHeatmapNodes(heatmap.Nodes),
		})
	}

	return output
}

// mapHeatmapNodes converts application heatmap nodes to transport output.
func mapHeatmapNodes(nodes []application.HeatmapNode) []HeatmapNode {
	output := make([]HeatmapNode, 0, len(nodes))
	for _, node := range nodes {
		output = append(output, HeatmapNode{
			Date:  time.Time(node.Date).Format(time.DateOnly),
			Value: node.Value,
		})
	}

	return output
}
