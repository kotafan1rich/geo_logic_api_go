package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

func ToType(infraType *InfrastructureType) *model.InfrastructureType {
	return &model.InfrastructureType{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func ToTypeModel(infraType *model.InfrastructureType) InfrastructureType {
	return InfrastructureType{
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func ToObject(infraObj *InfrastructureObject) *model.InfrastructureObject {
	return &model.InfrastructureObject{
		ID:       infraObj.ID,
		GeoPoint: model.GeoPoint(infraObj.Location),
		Address:  infraObj.Address,
		Name:     infraObj.Name,
		Type:     *ToType(&infraObj.Type),
	}
}

func ToObjectModel(infraObj *model.InfrastructureObject) *InfrastructureObject {
	return &InfrastructureObject{
		Location: database.DBGeoPoint(infraObj.GeoPoint),
		Address:  infraObj.Address,
		Name:     infraObj.Name,
		TypeID:   infraObj.Type.ID,
		Type:     ToTypeModel(&infraObj.Type),
	}
}

func ToObjectSlice(objs []InfrastructureObject) []model.InfrastructureObject {
	result := make([]model.InfrastructureObject, 0, len(objs))
	for _, infraObj := range objs {
		result = append(result, *ToObject(&infraObj))
	}
	return result
}
