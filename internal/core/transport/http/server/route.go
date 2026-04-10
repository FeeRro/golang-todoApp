package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/FeeRro/golang-todoApp/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string // /user/{id}
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
