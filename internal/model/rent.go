package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

type Rent struct {
	ID      uint64
	Lat     float64
	Long    float64
	Address string
	Info    string
}

func NewRent(lat, long float64, address, info string) (*Rent, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.ErrInvalidLat
	}
	if long < -180 || long > 180 {
		return nil, errors.ErrInvalidLong
	}
	return &Rent{
		Lat:     lat,
		Long:    long,
		Address: address,
		Info:    info,
	}, nil
}
