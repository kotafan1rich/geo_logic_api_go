package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
)

func ErrorHandlerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := errors.AsType[*apperror.AppError](err); ok {
				log.Warn("application error", "error", appErr)
				c.JSON(appErr.Status, appErr)
				c.Errors = nil
				return
			}
			if appErr, ok := errors.AsType[*gin.Error](err); ok {
				if appErr.Type == gin.ErrorTypeBind {
					log.Warn(
						"validation error",
						slog.String("error", appErr.Error()),
					)
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   "Validation failed",
						"code":    "VALIDATION_ERROR",
						"details": appErr.Error(),
					})
				} else {
					log.LogError(c.Request.Context(), appErr, "internal error")
					c.JSON(apperror.ErrInternal.Status, apperror.ErrInternal)
				}
				c.Errors = nil
				return
			}
			log.LogError(c.Request.Context(), err, "internal error")
			c.JSON(apperror.ErrInternal.Status, apperror.ErrInternal)

			c.Errors = nil
		}
	}
}
