package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/middleware"
)

type HttpHandler struct {
	log              logger.Logger
	userHandler      handler.UserHandler
	rentHandler      handler.RentHandler
	eventHandler     handler.EventHandler
	infraTypeHandler handler.InfraTypeHandler
	infraHandler     handler.InfraHandler
}

func NewHttpHandler(log logger.Logger, userHandler handler.UserHandler, rentHandler handler.RentHandler, eventHandler handler.EventHandler, infraTypeHandler handler.InfraTypeHandler, infraHandler handler.InfraHandler) *HttpHandler {
	return &HttpHandler{
		log:              log,
		userHandler:      userHandler,
		rentHandler:      rentHandler,
		eventHandler:     eventHandler,
		infraTypeHandler: infraTypeHandler,
		infraHandler:     infraHandler,
	}
}

func (h *HttpHandler) Routes() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		middleware.GinLoggerMiddleware(h.log),
		middleware.ErrorHandlerMiddleware(h.log),
	)

	api := router.Group("/api")

	handler.RegisterRoutes(api, h.userHandler, h.rentHandler, h.eventHandler, h.infraTypeHandler, h.infraHandler)

	return router
}
