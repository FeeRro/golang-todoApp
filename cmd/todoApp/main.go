package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/FeeRro/golang-todoApp/internal/core/config"
	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_pgx_pool "github.com/FeeRro/golang-todoApp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/FeeRro/golang-todoApp/internal/core/transport/http/middleware"
	core_http_server "github.com/FeeRro/golang-todoApp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/FeeRro/golang-todoApp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/FeeRro/golang-todoApp/internal/features/statistics/service"
	statistics_transport_http "github.com/FeeRro/golang-todoApp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/FeeRro/golang-todoApp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/FeeRro/golang-todoApp/internal/features/tasks/service"
	tasks_transport_http "github.com/FeeRro/golang-todoApp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/FeeRro/golang-todoApp/internal/features/users/repository/postgres"
	users_service "github.com/FeeRro/golang-todoApp/internal/features/users/service"
	users_transport_http "github.com/FeeRro/golang-todoApp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init aplication logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

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

	logger.Debug("initialized features", zap.String("features", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	taskService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHttp := tasks_transport_http.NewTasksHTTPHandler(taskService)

	logger.Debug("initialezed features", zap.String("features", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHttp := statistics_transport_http.NewStatisticsHTTPHandler(statisticService)

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
	apiVersionRouterV1.RegisterRoute(tasksTransportHttp.Routes()...)
	apiVersionRouterV1.RegisterRoute(statisticsTransportHttp.Routes()...)

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
