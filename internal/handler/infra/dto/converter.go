package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToTypeResponse(infraType *model.InfraType) *TypeResponse {
	return &TypeResponse{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}

func ToInfraResponse(infra *model.InfraObject) *InfraResponse {
	return &InfraResponse{
		ID:      infra.ID,
		Lat:     infra.GeoPoint.Lat,
		Long:    infra.GeoPoint.Long,
		Address: infra.Address,
		Name:    infra.Name,
		Type:    infra.Type.Name,
	}
}

func ToInfraResponseSlice(infras []model.InfraObject) []*InfraResponse {
	result := make([]*InfraResponse, 0, len(infras))
	for _, infra := range infras {
		result = append(result, ToInfraResponse(&infra))
	}
	return result
}
