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
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Available(c *gin.Context)
}

type EventHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Near(c *gin.Context)
}

type InfrastructureTypeHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func RegisterRoutes(router *gin.RouterGroup, userHandler UserHandler, rentHandler RentHandler, eventHandler EventHandler, typeHandler InfrastructureTypeHandler) {
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
		rent.PATCH("/:id", rentHandler.Update)
		rent.DELETE("/:id", rentHandler.Delete)
		rent.GET("/available", rentHandler.Available)
	}

	event := router.Group("/events")
	{
		event.POST("/", eventHandler.Create)
		event.GET("/:id", eventHandler.GetById)
		event.PATCH("/:id", eventHandler.Update)
		event.DELETE("/:id", eventHandler.Delete)
		event.GET("/near", eventHandler.Near)
	}

	infrastructure := router.Group("/infrastructure")
	{
		infraType := infrastructure.Group("/types")
		{
			infraType.POST("/", typeHandler.Create)
			infraType.GET("/:id", typeHandler.GetById)
			infraType.PATCH("/:id", typeHandler.Update)
			infraType.DELETE("/:id", typeHandler.Delete)
		}
	}

}
