package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToTypeResponse(infraType *model.InfraType) TypeResponse {
	return TypeResponse{
		ID:        infraType.ID,
		Slug:      infraType.Slug,
		Name:      infraType.Name,
		Weight:    infraType.Weight,
		MaxRadius: infraType.MaxRadius,
	}
}
