package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
)

// DeleteTask 	godoc
// @Summary 	Удалить задачу
// @Description Удаляет задачу по ID
// @Tags 		tasks
// @Param 		id path int true 								"ID задачи"
// @Success 	204 											"Задача успешно удалёна"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse 	"Task not found"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	task, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task id from path",
		)
		return
	}

	if err := h.tasksService.DeleteTask(ctx, task); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failet to delete task from repository",
		)
		return
	}

	responseHandler.NoContentResponse()

}
