package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type RentHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	// Available(c *gin.Context)
	// Delete(c *gin.Context)
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
		// rent.GET("/available", rentHandler.Available)
		// rent.DELETE("/:id", rentHandler.Delete)

	}
}
