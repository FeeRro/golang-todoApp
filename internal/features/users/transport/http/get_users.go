package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
	core_http_utils "github.com/FeeRro/golang-todoApp/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

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
	limit, err := core_http_utils.GetIntQueryParam(*r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(*r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
