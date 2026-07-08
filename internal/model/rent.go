package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

const (
	maxLatitude  = 90.0
	minLatitude  = -90.0
	maxLongitude = 180.0
	minLongitude = -180.0
)

type Rent struct {
	ID       uint64
	Geopoint *GeoPoint
	Address  string
	Info     *string
}

func NewRent(geopoint *GeoPoint, address string, info *string) *Rent {
	return &Rent{
		Geopoint: geopoint,
		Address:  address,
		Info:     info,
	}
}

func (r *Rent) UpdateAddress(address string) error {
	if address == "" {
		return errors.ErrInvalidAddress
	}
	r.Address = address
	return nil
}

func (r *Rent) UpdateInfo(info string) {
	r.Info = &info
}
