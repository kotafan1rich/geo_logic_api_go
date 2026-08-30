package trackedlocation

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	gendto "github.com/kotafan1rich/geo_logic_api_go/internal/handler/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/tracked_location/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type TrackedLocationService interface {
	Create(ctx context.Context, userID uint64, lat, long float64) (*model.TrackedLocation, error)
	GetByID(ctx context.Context, id uint64) (*model.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uint64) ([]model.TrackedLocation, error)
	Update(ctx context.Context, id uint64, userID *uint64, lat, long *float64) (*model.TrackedLocation, error)
	Delete(ctx context.Context, id uint64) error
}

type trackedLocationHandler struct {
	service TrackedLocationService
}

func NewHandler(service TrackedLocationService) *trackedLocationHandler {
	return &trackedLocationHandler{service: service}
}

func (h *trackedLocationHandler) Create(c *gin.Context) {
	var req dto.CreateTrackedLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	location, err := h.service.Create(c.Request.Context(), req.UserID, req.Lat, req.Long)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToTrackedLocationResponse(location))
}

func (h *trackedLocationHandler) GetByID(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	location, err := h.service.GetByID(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTrackedLocationResponse(location))
}

func (h *trackedLocationHandler) GetByUserID(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	locations, err := h.service.GetByUserID(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTrackedLocationResponseSlice(locations))
}

func (h *trackedLocationHandler) Update(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	var reqBody dto.UpdateTrackedLocationRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	location, err := h.service.Update(
		c.Request.Context(),
		reqUri.ID,
		reqBody.UserID,
		reqBody.Lat,
		reqBody.Long,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ToTrackedLocationResponse(location))
}

func (h *trackedLocationHandler) Delete(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.service.Delete(c.Request.Context(), reqUri.ID); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
