package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
)

func GinLoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		requestLogger := log.WithRequest(
			status,
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.URL.RawQuery,
			c.ClientIP(),
			c.Request.UserAgent(),
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
