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

func ToRentSlice(rents []Rent) []model.Rent {
	result := make([]model.Rent, 0, len(rents))
	for _, rent := range rents {
		result = append(result, *ToRent(&rent))
	}
	return result
}

func ToRentModel(rent *model.Rent) *Rent {
	res := &Rent{
		Location: DBGeoPoint(*rent.Geopoint),
		Address:  rent.Address,
		Info:     rent.Info,
	}
	res.ID = rent.ID
	return res
}
