package rent

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	gendto "github.com/kotafan1rich/geo_logic_api_go/internal/handler/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type Logger interface {
	Error(msg string, args ...any)
}

type rentHandler struct {
	service service.RentService
	log  Logger
}

func NewHandler(service service.RentService, log Logger) *rentHandler {
	return &rentHandler{service: service, log: log}
}

func (h *rentHandler) Create(c *gin.Context) {
	var req dto.CreateRentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	rent, err := h.service.Create(c.Request.Context(), req.Lat, req.Long, req.Address, req.Info)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToRentResponse(rent))
}

func (h *rentHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.Error(errors.ValidationError("id must be greater then 0")).SetType(gin.ErrorTypeBind)
		return
	}

	rent, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToRentResponse(rent))
}

func (h *rentHandler) Update(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	var reqBody dto.UpdateRentRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errDetails := handler.ParseValidationError(err)
		c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	updatedRent, err := h.service.Update(
		c.Request.Context(),
		reqUri.ID,
		reqBody.Lat,
		reqBody.Long,
		reqBody.Address,
		reqBody.Info,
	)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToRentResponse(updatedRent))
}

func (h *rentHandler) Delete(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	err := h.service.Delete(c.Request.Context(), reqUri.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *rentHandler) Available(c *gin.Context) {
	var req dto.AvailableRentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	geopoint, err := model.NewGeoPoint(req.Lat, req.Long)
	if err != nil {
		c.Error(errors.ValidationError(err.Error())).SetType(gin.ErrorTypeBind)
		return
	}
	results, err := h.service.Available(c.Request.Context(), geopoint, req.Radius)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToRentResponseSlice(results))
}
