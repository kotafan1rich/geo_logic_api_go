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

func (r *Rent) Update(lat, long *float64, address *string, info *string) error {
	if lat != nil {
		err := r.Geopoint.UpdateLat(*lat)
		if err != nil {
			return err
		}
	}

	if long != nil {
		err := r.Geopoint.UpdateLong(*long)
		if err != nil {
			return err
		}
	}

	if address != nil {
		err := r.UpdateAddress(*address)
		if err != nil {
			return err
		}
	}

	if info != nil {
		r.UpdateInfo(*info)
	}
	return nil
}
