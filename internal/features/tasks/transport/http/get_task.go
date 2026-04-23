package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask 		godoc
// @Summary 	Получить задачу
// @Description Получать задачу по ID
// @Tags 		tasks
// @Produce 	json
// @Param 		id path int true 								"ID задачи"
// @Success 	200 {object} GetTaskResponse 					"Задача найдена"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse 	"Task not found"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"Failed to get task by id",
		)
		return
	}

	task, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task from repository",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(task))

	responseHandler.JsonResponse(response, http.StatusOK)
}
