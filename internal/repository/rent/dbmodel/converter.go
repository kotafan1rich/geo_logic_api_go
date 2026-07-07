package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

func ToRent(rent *Rent) *model.Rent {
	return &model.Rent{
		ID:       rent.ID,
		Geopoint: (*model.GeoPoint)(&rent.Location),
		Address:  rent.Address,
		Info:     rent.Info,
	}
}

func ToRentModel(rent *model.Rent) *Rent {
	return &Rent{
		Location: DBGeoPoint(*rent.Geopoint),
		Address:  rent.Address,
		Info:     rent.Info,
	}
}
