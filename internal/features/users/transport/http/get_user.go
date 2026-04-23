package users_transport_http

import (
	"net/http"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser 		godoc
// @Summary 	Получить пользователя
// @Description Получать пользователя по его ID
// @Tags 		users
// @Produce 	json
// @Param 		id path int true 								"ID пользователя"
// @Success 	200 {object} GetUserResponse 					"Пользователь найден"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse 	"User not found"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/users/{id} [get]
func (h UsersHttpHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")
		return
	}

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JsonResponse(response, http.StatusOK)
}
