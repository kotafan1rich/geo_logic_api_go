package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	gendto "github.com/kotafan1rich/geo_logic_api_go/internal/handler/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/user/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type userHandler struct {
	service service.UserService
}

func NewHandler(service service.UserService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.service.Create(c.Request.Context(), req.TgID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (h *userHandler) GetByID(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (h *userHandler) Update(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	var reqBody dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	user, err := h.service.Update(c.Request.Context(), reqUri.ID, reqBody.TgID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (h *userHandler) Delete(c *gin.Context) {
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
