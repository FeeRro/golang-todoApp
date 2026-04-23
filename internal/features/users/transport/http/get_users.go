package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers 	godoc
// @Summary 	Список пользователей
// @Description Получать список пользователей с опциональной пагинацией (limit, offset)
// @Tags 		users
// @Produce 	json
// @Param 		limit query int false 							"limit - количество пользователей для получения"
// @Param 		offset query int false 							"offset - смещение от начала списка пользователей"
// @Success 	200 {object} GetUsersResponse 					"Пользователи найдены"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/users [get]
func (h *UsersHttpHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParam(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed 'limit'/'offset' query param")
		return
	}

	userDomains, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParam(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParam  = "limit"
		offsetQueryParam = "offset"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParam)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
