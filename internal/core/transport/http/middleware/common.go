package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	core_logger "github.com/FeeRro/golang-todoApp/internal/core/logger"
	core_http_response "github.com/FeeRro/golang-todoApp/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Генерация уникального идентификатора для каждого запроса
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString() // Функция для генерации уникального ID
			}

			r.Header.Set("X-Request-ID", requestID)   // Установка ID в заголовок запроса
			w.Header().Set("X-Request-ID", requestID) // Установка ID в заголовок ответа

			next.ServeHTTP(w, r) // Вызов следующего обработчика
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if pan := recover(); pan != nil {
					responseHandler.PanicResponse(
						pan,
						"Паника при обработке Http запроса",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCodeOrPanic()),
				zap.Duration("Latency", time.Since(before)),
			)
		})
	}
}
