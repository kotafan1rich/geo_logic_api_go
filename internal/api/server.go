package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
)

type Logger interface {
	GinLoggerMiddleware() gin.HandlerFunc
}

type Handler interface {
	Routes() *gin.Engine
}

type mainHandler struct {
	log         Logger
	userHandler handler.UserHandler
}

func NewMainHandler(log Logger, userHandler handler.UserHandler) Handler {
	return &mainHandler{log: log, userHandler: userHandler}
}

func (h *mainHandler) Routes() *gin.Engine {
	router := gin.New()

	router.Use(
		h.log.GinLoggerMiddleware(),
		gin.Recovery(),
	)

	api := router.Group("/api")

	handler.RegisterRoutes(api, h.userHandler)

	return router
}
