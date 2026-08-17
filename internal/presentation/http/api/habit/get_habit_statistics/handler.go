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
		MeasurableSeries:  mapSeries(result.MeasurableSeries),
		TimeSeries:        mapSeries(result.TimeSeries),
		CompletableSeries: mapSeries(result.CompletableSeries),
	}
}

// mapSeries converts application series to transport output.
func mapSeries(series []application.Series) []Series {
	output := make([]Series, 0, len(series))
	for _, item := range series {
		output = append(output, Series{
			HabitID:  item.HabitID.String(),
			Nodes:    mapNodes(item.Nodes),
			MinValue: item.MinValue,
			MaxValue: item.MaxValue,
		})
	}

	return output
}

// mapNodes converts application nodes to transport output.
func mapNodes(nodes []application.Node) []Node {
	output := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		output = append(output, Node{
			Date:  time.Time(node.Date).Format(time.DateOnly),
			Value: node.Value,
		})
	}

	return output
}
