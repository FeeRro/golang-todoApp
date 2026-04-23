package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/FeeRro/golang-todoApp/internal/core/domain"
	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_request "github.com/FeeRro/golang-todoApp/internal/core/transport/http/request"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
	core_http_types "github.com/FeeRro/golang-todoApp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"Ivan Ivanov"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79998887766"`
}

// PatchUser 		godoc
// @Summary 	Изменение пользователя
// @Description Изменение информации о уже существующем в системе пользователе по его ID
// @Description ### Логика обновления полей (Three-state logec):
// @Description 1. **Поле не передано**: 'phone_number' игнорируется, значение в БД не изменяется
// @Description 2. **Поле явно передано**: '"phone_number": "+79998887766"' - устанавливает новый номер в БД
// @Description 3. **Передан null**: '"phone_number": null' - удаляет номер телефона в БД (set to NULL)
// @Description Ограничения: 'full_name' не может быть выставлен как null
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		id path int true 								"ID пользователя"
// @Param 		request body PatchUserRequest true 				"PatchUser тело запроса"
// @Success 	200 {object} GetUserResponse 					"Успешное изменение пользоватнля"
// @Failure 	400 {object} core_http_response.ErrorResponse 	"Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse 	"User not found"
// @Failure 	500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/users/{id} [patch]
func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'fullName' can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("'fullName' must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("'PhoneNumber' must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("'PhoneNumber' must startswith '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHttpHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failet to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userPatch := UserPatchFromRequest(request)

	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.Todomain(),
		request.PhoneNumber.Todomain(),
	)
}
