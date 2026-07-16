package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	gendto "github.com/kotafan1rich/geo_logic_api_go/internal/handler/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/event/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type eventHandler struct {
	service service.EventService
}

func NewHandler(service service.EventService) *eventHandler {
	return &eventHandler{service: service}
}

func (h *eventHandler) Create(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	event, err := h.service.Create(c.Request.Context(), req.Lat, req.Long, req.Date, req.Info)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToEventResponse(event))
}

func (h *eventHandler) GetById(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	event, err := h.service.GetById(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToEventResponse(event))
}

func (h *eventHandler) Update(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	var reqBody dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	updatedEvent, err := h.service.Update(
		c.Request.Context(),
		reqUri.ID,
		reqBody.Lat,
		reqBody.Long,
		reqBody.Date,
		reqBody.Info,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToEventResponse(updatedEvent))
}

func (h *eventHandler) Delete(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	err := h.service.Delete(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *eventHandler) Near(c *gin.Context) {
	var req dto.AvailableEventRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	geopoint, err := model.NewGeoPoint(req.Lat, req.Long)
	if err != nil {
		_ = c.Error(errors.ValidationError(err.Error())).SetType(gin.ErrorTypeBind)
		return
	}
	results, err := h.service.Near(c.Request.Context(), geopoint, req.Radius)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToEventResponseSlice(results))
}
