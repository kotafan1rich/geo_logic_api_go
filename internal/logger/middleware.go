package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (l *logger) GinLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		requestLogger := l.WithRequest(
			status,
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Request.URL.RawQuery,
			ctx.ClientIP(),
			ctx.Request.UserAgent(),
			latency,
		)

		switch {
		case status >= 500:
			requestLogger.Error("Server error")
		case status >= 400:
			requestLogger.Warn("Client error")
		default:
			requestLogger.Info("Request complete")
		}
	}
}
