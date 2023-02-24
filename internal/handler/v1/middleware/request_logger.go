package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (mw *Manager) RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()

		mw.logger.Infow(
			"middleware logger",
			"method", req.Method,
			"URI", req.URL,
			"status", status,
			"size", size,
			"time", s)
		return err
	}
}
