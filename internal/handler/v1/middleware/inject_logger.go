package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/malkev1ch/apod/pkg/logger"
)

func (mw *Manager) InjectLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return next(logger.ContextEchoWithLogger(ctx, mw.logger))
	}
}
