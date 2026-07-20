package infra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	gendto "github.com/kotafan1rich/geo_logic_api_go/internal/handler/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/infra/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type infraHandler struct {
	service service.InfraService
}

func NewInfraHandler(service service.InfraService) *infraHandler {
	return &infraHandler{service: service}
}

func (h *infraHandler) Create(c *gin.Context) {
	var req dto.CreateInfraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	infra, err := h.service.Create(c.Request.Context(), req.Lat, req.Long, req.Address, req.Name, req.TypeId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToInfraResponse(infra))
}

func (h *infraHandler) GetByID(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	infra, err := h.service.GetByID(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToInfraResponse(infra))
}

func (h *infraHandler) Update(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	var reqBody dto.UpdateInfraRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	updatedInfra, err := h.service.Update(
		c.Request.Context(),
		reqUri.ID,
		reqBody.Lat,
		reqBody.Long,
		reqBody.Address,
		reqBody.Name,
		reqBody.TypeId,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToInfraResponse(updatedInfra))
}

func (h *infraHandler) Delete(c *gin.Context) {
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

func (h *infraHandler) Near(c *gin.Context) {
	var req dto.NearEventRequest
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
	results, err := h.service.Near(c.Request.Context(), geopoint)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToInfraResponseSlice(results))
}
