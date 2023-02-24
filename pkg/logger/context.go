package logger

import (
	"context"
	"github.com/labstack/echo/v4"
)

// loggerKey points to the value in the context where the logger is stored.
const loggerKey = "logger"

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// LoggerFromContext returns logger from context.
func LoggerFromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey).(*appLogger); ok {
		return l
	}
	return defLogger
}

// ContextEchoPropagateLogger adds logger to underlay context.
func ContextEchoPropagateLogger(ctx echo.Context) context.Context {
	return ContextWithLogger(ctx.Request().Context(), LoggerFromContextEcho(ctx))
}

// ContextEchoWithLogger adds logger to echo context.
func ContextEchoWithLogger(ctx echo.Context, l Logger) echo.Context {
	ctx.Set(string(loggerKey), l)
	return ctx
}

// LoggerFromContextEcho returns logger from context.
func LoggerFromContextEcho(ctx echo.Context) Logger {
	if l, ok := ctx.Get(loggerKey).(*appLogger); ok {
		return l
	}
	return defLogger
}