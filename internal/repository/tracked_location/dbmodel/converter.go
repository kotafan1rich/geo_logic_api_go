package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

func ToDomain(location *TrackedLocation) *model.TrackedLocation {
	return &model.TrackedLocation{
		ID:       location.ID,
		UserID:   location.UserID,
		GeoPoint: model.GeoPoint(location.Location),
	}
}

func ToDomainSlice(rents []TrackedLocation) []model.TrackedLocation {
	result := make([]model.TrackedLocation, 0, len(rents))
	for _, location := range rents {
		result = append(result, *ToDomain(&location))
	}
	return result
}

func ToModel(location *model.TrackedLocation) *TrackedLocation {
	res := &TrackedLocation{
		UserID:   location.UserID,
		Location: database.DBGeoPoint(location.GeoPoint),
	}
	res.ID = location.ID
	return res
}
