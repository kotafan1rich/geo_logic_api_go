package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
)

func ErrorHandlerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case *errors.AppError:
				log.Warn("application error", "error", e)
				c.JSON(e.Status, e)
			case *gin.Error:
				if e.Type == gin.ErrorTypeBind {
					log.Warn("validation error", "error", e)
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   "Validation failed",
						"code":    "VALIDATION_ERROR",
						"details": e.Error(),
					})
				} else {
					log.LogError(c.Request.Context(), e, "internal error")
					c.JSON(errors.ErrInternal.Status, errors.ErrInternal)
				}
			default:
				log.LogError(c.Request.Context(), e, "internal error")
				c.JSON(errors.ErrInternal.Status, errors.ErrInternal)

			}
			c.Errors = nil
		}
	}
}
