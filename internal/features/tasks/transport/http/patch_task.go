package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/FeeRro/golang-todoApp/internal/core/domain"
	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
	core_http_types "github.com/FeeRro/golang-todoApp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"        swaggertype:"string"  example:"Выбросить мусор"`
	Description core_http_types.Nullable[string] `json:"description"  swaggertype:"string"  example:"Выбросить пока не пришла мама"`
	Completed   core_http_types.Nullable[bool]   `json:"is_completed" swaggertype:"boolean" example:"true"`
}

// PatchTask 	godoc
// @Summary 	Изменение задачи
// @Description Изменение информации о уже существующем в системе задачи по ID
// @Description ### Логика обновления полей (Three-state logec):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не изменяется
// @Description 2. **Поле явно передано**: `"description": "Сделать матан до 18:00" - устанавливает новое описание в БД
// @Description 3. **Передан null**: `"description": null` - удаляет описание в БД (set to NULL)
// @Description Ограничения: 'title' и `is_completed` не может быть выставлен как null
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		id path int true 								"ID задачи"
// @Param 		request body PatchTaskRequest true 				"PatchTask тело запроса"
// @Success 	200 {object} GetTaskResponse 					"Успешное изменение задачи"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse 	"Task not found"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/tasks/{id} [patch]
func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'TITLE' can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("'Completed' can't be NULL")
		}
	}

	return nil
}

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JsonResponse(response, http.StatusOK)

}

type PatchTaskResponse TaskDTOResponse

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.Todomain(),
		request.Description.Todomain(),
		request.Completed.Todomain(),
	)
}
