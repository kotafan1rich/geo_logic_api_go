package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToRentResponse(rent *model.Rent) *RentResponse {
	return &RentResponse{
		ID:      rent.ID,
		Lat:     rent.Geopoint.Lat,
		Long:    rent.Geopoint.Long,
		Address: rent.Address,
		Info:    rent.Info,
	}
}
