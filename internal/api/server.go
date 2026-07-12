package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/middleware"
)

type Handler interface {
	Routes() *gin.Engine
}

type mainHandler struct {
	log         logger.Logger
	userHandler handler.UserHandler
	rentHandler handler.RentHandler
}

func NewMainHandler(log logger.Logger, userHandler handler.UserHandler, rentHandler handler.RentHandler) Handler {
	return &mainHandler{log: log, userHandler: userHandler, rentHandler: rentHandler}
}

func (h *mainHandler) Routes() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		middleware.GinLoggerMiddleware(h.log),
		middleware.ErrorHandlerMiddleware(h.log),
	)

	api := router.Group("/api")

	handler.RegisterRoutes(api, h.userHandler, h.rentHandler)

	return router
}
