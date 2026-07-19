package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

func ToType(infraType *InfraType) *model.InfraType {
	return &model.InfraType{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func ToTypeModel(infraType *model.InfraType) InfraType {
	return InfraType{
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func ToObject(infraObj *InfraObject) *model.InfraObject {
	return &model.InfraObject{
		ID:       infraObj.ID,
		GeoPoint: model.GeoPoint(infraObj.Location),
		Address:  infraObj.Address,
		Name:     infraObj.Name,
		Type:     *ToType(&infraObj.Type),
	}
}

func ToObjectModel(infraObj *model.InfraObject) *InfraObject {
	return &InfraObject{
		Location: database.DBGeoPoint(infraObj.GeoPoint),
		Address:  infraObj.Address,
		Name:     infraObj.Name,
		TypeID:   infraObj.Type.ID,
		Type:     ToTypeModel(&infraObj.Type),
	}
}

func ToObjectSlice(objs []InfraObject) []model.InfraObject {
	result := make([]model.InfraObject, 0, len(objs))
	for _, infraObj := range objs {
		result = append(result, *ToObject(&infraObj))
	}
	return result
}
