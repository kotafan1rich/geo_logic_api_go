package handler

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Create(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func RegisterRoutes(router *gin.RouterGroup, userHandler UserHandler) {
	user := router.Group("/user")
	{
		user.POST("/create", userHandler.Create)
		user.GET("/get_by_id/:id", userHandler.GetById)
		user.PUT("/update", userHandler.Update)
		user.DELETE("/delete/:id", userHandler.Delete)
	}
}
