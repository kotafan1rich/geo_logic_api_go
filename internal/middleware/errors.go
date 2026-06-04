package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			switch e := err.(type) {
			case *errors.AppError:
				ctx.JSON(e.Status, e)
			case *gin.Error:
				if e.Type == gin.ErrorTypeBind {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   "Validation failed",
						"code":    "VALIDATION_ERROR",
						"details": e.Error(),
					})
				} else {
					ctx.JSON(errors.ErrInternal.Status, errors.ErrInternal)
				}
			default:
				ctx.JSON(errors.ErrInternal.Status, errors.ErrInternal)

			}
			ctx.Errors = nil
		}
	}
}
