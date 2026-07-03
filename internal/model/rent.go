package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

const (
	maxLatitude  = 90.0
	minLatitude  = -90.0
	maxLongitude = 180.0
	minLongitude = -180.0
)

type Rent struct {
	ID      uint64
	Lat     float64
	Long    float64
	Address string
	Info    string
}

func NewRent(lat, long float64, address, info string) (*Rent, error) {
	if lat < minLatitude || lat > maxLatitude {
		return nil, errors.ErrInvalidLat
	}
	if long < minLongitude || long > maxLongitude {
		return nil, errors.ErrInvalidLong
	}
	return &Rent{
		Lat:     lat,
		Long:    long,
		Address: address,
		Info:    info,
	}, nil
}
