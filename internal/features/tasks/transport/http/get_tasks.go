package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks 	godoc
// @Summary 	Список задач
// @Description Получать список задач с опциональной пагинацией и/или фильтрацией по ID (userID, limit, offset)
// @Tags 		tasks
// @Produce 	json
// @param       user_id query int false 						"user_id - ID пользователя, задачи которого нужно получить"
// @Param 		limit query int false 							"limit - количество пользователей для получения"
// @Param 		offset query int false 							"offset - смещение от начала списка пользователей"
// @Success 	200 {object} GetTasksResponse 					"Список задач"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/tasks [get]
func (h TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userId/limit/offset query params",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(tasksDTOFromDomains(tasksDomains))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParam(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParam = "user_id"
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
