package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Create(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type RentHandler interface {
	Create(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Update(ctx *gin.Context)
	// Available(ctx *gin.Context)
	// Delete(ctx *gin.Context)
}

func RegisterRoutes(router *gin.RouterGroup, userHandler UserHandler, rentHandler RentHandler) {
	user := router.Group("/users")
	{
		user.POST("/", userHandler.Create)
		user.GET("/:id", userHandler.GetById)
		user.PATCH("/:id", userHandler.Update)
		user.DELETE("/:id", userHandler.Delete)
	}

	rent := router.Group("/rents")
	{
		rent.POST("/", rentHandler.Create)
		rent.GET("/:id", rentHandler.GetById)
		user.PATCH("/:id", rentHandler.Update)
		// rent.GET("/available", rentHandler.Available)
		// rent.DELETE("/:id", rentHandler.Delete)

	}
}
