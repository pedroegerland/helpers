package http

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

// Chain method chains the middlewares to execute before handler
// ref: https://gist.github.com/husobee/fd23681261a39699ee37.
func Chain(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](Chain(h, m[1:]...))
}

func Server(httpServer *http.Server) fx.Hook {
	return fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				zap.L().Info(
					"http_server_up",
					zap.String("description", "up and running api server"),
					zap.String("address", httpServer.Addr),
				)
				if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
					zap.L().Error("http_server_down", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	}
}

func NewTracerHTTPMiddleware(ignorePaths ...string) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			path := request.URL.Path
			for _, ignorePath := range ignorePaths {
				if strings.Contains(path, ignorePath) {
					handler.ServeHTTP(writer, request)
					return
				}
			}

			wrapWriter := &WrapResponseWriter{ResponseWriter: writer}
			ctx := request.Context()

			request = request.WithContext(ctx)
			handler.ServeHTTP(wrapWriter, request)
			statusCode := wrapWriter.statusCode

			if statusCode >= http.StatusBadRequest {
				zap.L().Error(fmt.Sprintf("response code: %d", statusCode))
			}
		})
	}
}
