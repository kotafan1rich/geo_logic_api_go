package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

func ToRent(rent *Rent) *model.Rent {
	return &model.Rent{
		ID:      rent.ID,
		Lat:     rent.Location.Lat,
		Long:    rent.Location.Long,
		Address: rent.Address,
		Info:    rent.Info,
	}
}

func ToRentModel(rent *model.Rent) *Rent {
	return &Rent{
		Location: DBGeoPoint{Lat: rent.Lat, Long: rent.Long},
		Address: rent.Address,
		Info:    rent.Info,
	}
}
