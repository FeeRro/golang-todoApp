package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/FeeRro/golang-todoApp/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	APIVersionV1 = APIVersion("v1")
	APIVersionV2 = APIVersion("v2")
	APIVersionV3 = APIVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(apiVersion APIVersion, middleware ...core_http_middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoute(routes ...Route) {
	for _, route := range routes {
		// "Get /tasks"
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
