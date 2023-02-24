package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/malkev1ch/apod/gen/v1"
	berror "github.com/malkev1ch/apod/pkg/errors"
)

func (mw *Manager) BusinessError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		if err != nil {
			mw.logger.Error(err)
			var e *berror.Error
			if errors.As(err, &e) {
				message, httpCode := parseHTTPCode(e.Message, e.Code)
				return ctx.JSON(
					httpCode,
					gen.ErrorResponse{
						Message: message,
					},
				)
			}

			var eEcho *echo.HTTPError
			if errors.As(err, &eEcho) {
				return ctx.JSON(
					eEcho.Code,
					gen.ErrorResponse{
						Message: fmt.Sprintf("%v", eEcho.Message),
					},
				)
			}

			return err
		}

		return nil
	}
}

func parseHTTPCode(message string, code uint32) (string, int) {
	switch code {
	case berror.NasaPictureInvalidFormatErrCode:
		return message, http.StatusInternalServerError
	default:
		return "Internal Server Error", http.StatusInternalServerError
	}
}
