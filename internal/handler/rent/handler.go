package rent

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type rentHandler struct {
	service service.RentService
}

func NewHandler(service service.RentService) *rentHandler {
	return &rentHandler{service: service}
}

func (h *rentHandler) Create(ctx *gin.Context) {
	var req dto.CreateRentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		ctx.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	rent, err := h.service.Create(ctx.Request.Context(), req.Lat, req.Long, req.Address, req.Info)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToRentResponse(rent))
}

func (h *rentHandler) GetById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.Error(errors.ValidationError("id must be greater then 0")).SetType(gin.ErrorTypeBind)
		return
	}

	rent, err := h.service.GetById(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.ToRentResponse(rent))
}

func (h *rentHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.Error(errors.ValidationError("id must be greater then 0")).SetType(gin.ErrorTypeBind)
		return
	}
	var req dto.UpdateRentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		ctx.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	updatedRent, err := h.service.Update(ctx, id, req.Lat, req.Long, req.Address, req.Info)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.ToRentResponse(updatedRent))
}
