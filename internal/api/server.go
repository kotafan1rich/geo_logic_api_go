package api

import "github.com/gin-gonic/gin"


type Handler interface {
	Routes() *gin.Engine
}

type handler struct {
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) Routes() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
	)

	return router
}
