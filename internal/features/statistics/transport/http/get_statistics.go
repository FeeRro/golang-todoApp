package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FeeRro/golang-todoApp/internal/core/domain"
	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_create"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	queryParams, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, queryParams.userID, queryParams.from, queryParams.to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JsonResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

type queryParam struct {
	userID *int
	from   *time.Time
	to     *time.Time
}

func getUserIDFromToQueryParams(r *http.Request) (queryParam, error) {
	const (
		userIDQueryParam = "user_id"
		fromQueryParam   = "from"
		toQueryParam     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		return queryParam{}, fmt.Errorf(
			"get 'user_id' query param: %w",
			err,
		)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParam)
	if err != nil {
		return queryParam{}, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toQueryParam)
	if err != nil {
		return queryParam{}, fmt.Errorf("get 'to' query param: %w", err)
	}

	return queryParam{
		userID: userID,
		from:   from,
		to:     to,
	}, nil
}
