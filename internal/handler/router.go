package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type RentHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Available(c *gin.Context)
}

type EventHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Near(c *gin.Context)
}

type InfraTypeHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type InfraHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Near(c *gin.Context)
}

type TrackedLocationHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetByUserID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func RegisterRoutes(router *gin.RouterGroup, userHandler UserHandler, rentHandler RentHandler, eventHandler EventHandler, typeHandler InfraTypeHandler, infraHandler InfraHandler, trackedLocationHandler TrackedLocationHandler) {
	user := router.Group("/users")
	{
		user.POST("/", userHandler.Create)
		user.GET("/:id", userHandler.GetByID)
		user.PATCH("/:id", userHandler.Update)
		user.DELETE("/:id", userHandler.Delete)
	}

	rent := router.Group("/rents")
	{
		rent.POST("/", rentHandler.Create)
		rent.GET("/:id", rentHandler.GetByID)
		rent.PATCH("/:id", rentHandler.Update)
		rent.DELETE("/:id", rentHandler.Delete)
		rent.GET("/available", rentHandler.Available)
	}

	event := router.Group("/events")
	{
		event.POST("/", eventHandler.Create)
		event.GET("/:id", eventHandler.GetByID)
		event.PATCH("/:id", eventHandler.Update)
		event.DELETE("/:id", eventHandler.Delete)
		event.GET("/near", eventHandler.Near)
	}

	infra := router.Group("/infra")
	{
		infra.POST("/", infraHandler.Create)
		infra.GET("/:id", infraHandler.GetByID)
		infra.PATCH("/:id", infraHandler.Update)
		infra.DELETE("/:id", infraHandler.Delete)
		infra.GET("/near", infraHandler.Near)

		infraType := infra.Group("/types")
		{
			infraType.POST("/", typeHandler.Create)
			infraType.GET("/:id", typeHandler.GetByID)
			infraType.PATCH("/:id", typeHandler.Update)
			infraType.DELETE("/:id", typeHandler.Delete)
		}
	}

	trackedLocation := router.Group("/tracked-locations")
	{
		trackedLocation.POST("/", trackedLocationHandler.Create)
		trackedLocation.GET("/:id", trackedLocationHandler.GetByID)
		trackedLocation.GET("/user/:id", trackedLocationHandler.GetByUserID)
		trackedLocation.PATCH("/:id", trackedLocationHandler.Update)
		trackedLocation.DELETE("/:id", trackedLocationHandler.Delete)
	}
}
