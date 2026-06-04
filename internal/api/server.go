package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
)

type Handler interface {
	Routes() *gin.Engine
}

type mainHandler struct {
	userHandler handler.UserHandler
}

func NewMainHandler(userHandler handler.UserHandler) Handler {
	return &mainHandler{userHandler: userHandler}
}

func (h *mainHandler) Routes() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
	)

	api := router.Group("/api")
	{
		handler.RegisterRoutes(api, h.userHandler)
	}

	return router
}
