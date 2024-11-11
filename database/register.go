package databasse

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ReadinessFunc(chain Chain) func(echo.Context) error {
	return func(c echo.Context) error {
		if err := chain.Readiness(c.Request().Context()); err != nil {
			zap.L().Error("readiness_error", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func LivenessFunc(chain Chain) func(echo.Context) error {
	return func(c echo.Context) error {
		if err := chain.Liveness(c.Request().Context()); err != nil {
			zap.L().Error("liveness_error", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	}
}

func Register(e *echo.Echo, chain Chain) {
	e.GET("/readiness", ReadinessFunc(chain))
	e.GET("/liveness", LivenessFunc(chain))
}
