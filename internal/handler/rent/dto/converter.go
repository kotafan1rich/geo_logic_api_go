package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToRentResponse(rent *model.Rent) *RentResponse {
	return &RentResponse{
		Id:      rent.ID,
		Lat:     rent.Lat,
		Long:    rent.Long,
		Address: rent.Address,
		Info:    rent.Info,
	}
}
