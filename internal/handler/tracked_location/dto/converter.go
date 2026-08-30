package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToTrackedLocationResponse(location *model.TrackedLocation) *TrackedLocationResponse {
	return &TrackedLocationResponse{
		ID:     location.ID,
		UserID: location.UserID,
		Lat:    location.GeoPoint.Lat,
		Long:   location.GeoPoint.Long,
	}
}

func ToTrackedLocationResponseSlice(locations []model.TrackedLocation) []*TrackedLocationResponse {
	result := make([]*TrackedLocationResponse, 0, len(locations))
	for _, location := range locations {
		result = append(result, ToTrackedLocationResponse(&location))
	}
	return result
}
