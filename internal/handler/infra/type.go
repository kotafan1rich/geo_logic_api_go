package infra

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler"
	gendto "github.com/kotafan1rich/geo_logic_api_go/internal/handler/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/infra/dto"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type typeHandler struct {
	service service.InfraTypeService
}

func NewTypeHandler(service service.InfraTypeService) *typeHandler {
	return &typeHandler{service: service}
}

func (t *typeHandler) Create(c *gin.Context) {
	var req dto.CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	infraType, err := t.service.Create(c.Request.Context(), req.Slug, req.Name, req.Weight, req.MaxRadius)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToTypeResponse(infraType))
}

func (t *typeHandler) GetByID(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	infraType, err := t.service.GetByID(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToTypeResponse(infraType))
}

func (t *typeHandler) Update(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}
	var reqBody dto.UpdateTypeRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	infraType, err := t.service.Update(c.Request.Context(), reqUri.ID, reqBody.Slug, reqBody.Name, reqBody.Weight, reqBody.MaxRadius)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToTypeResponse(infraType))
}

func (t *typeHandler) Delete(c *gin.Context) {
	var reqUri gendto.IDUriRequest
	if err := c.ShouldBindUri(&reqUri); err != nil {
		errDetails := handler.ParseValidationError(err)
		_ = c.Error(errors.ValidationError(errDetails)).SetType(gin.ErrorTypeBind)
		return
	}

	err := t.service.Delete(c.Request.Context(), reqUri.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
