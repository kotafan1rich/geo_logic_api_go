package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/user/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type userHandler struct {
	service service.UserService
}

func NewHandler(service service.UserService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) Create(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		ctx.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.service.Create(ctx.Request.Context(), req.TgID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (h *userHandler) GetById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.Error(errors.ValidationError("id must be greater then 0")).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.service.GetById(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (h *userHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.Error(errors.ValidationError("id must be greater then 0")).SetType(gin.ErrorTypeBind)
		return
	}
	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		ctx.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.service.Update(ctx.Request.Context(), id, req.TgID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (h *userHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.Error(errors.ValidationError("id must be greater then 0")).SetType(gin.ErrorTypeBind)
		return
	}

	err = h.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
