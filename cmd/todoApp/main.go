package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_pgx_pool "github.com/FeeRro/golang-todoApp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/FeeRro/golang-todoApp/internal/core/transport/http/middleware"
	core_http_server "github.com/FeeRro/golang-todoApp/internal/core/transport/http/server"
	users_postgres_repository "github.com/FeeRro/golang-todoApp/internal/features/users/repository/postgres"
	users_service "github.com/FeeRro/golang-todoApp/internal/features/users/service"
	users_transport_http "github.com/FeeRro/golang-todoApp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init aplication logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initialized postgres connection pool")
	pool, err := core_pgx_pool.NewConnectionPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initialized features", zap.String("features", "users"))
	usersRepository := users_postgres_repository.NewUserRepository(pool)
	userService := users_service.NewUserService(usersRepository)

	usersTransportHttp := users_transport_http.NewUsersHttpHandler(userService)

	logger.Debug("initialized HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersionV1)
	apiVersionRouterV1.RegisterRoute(usersTransportHttp.Routes()...)

	// Пример работы мидлвари на apiVersionRouterV2

	// apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
	// 	core_http_server.APIVersionV2,
	// 	core_http_middleware.Dummy("тест мидлвари"),
	// )
	// apiVersionRouterV2.RegisterRoute(usersTransportHttp.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		// apiVersionRouterV2,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run Error", zap.Error(err))
	}
}
