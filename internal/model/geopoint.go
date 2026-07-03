package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

type GeoPoint struct {
	Lat  float64
	Long float64
}

func NewGeoPoint(lat, long float64) (*GeoPoint, error) {
	if lat < minLatitude || lat > maxLatitude {
		return nil, errors.ErrInvalidLat
	}
	if long < minLongitude || long > maxLongitude {
		return nil, errors.ErrInvalidLong
	}
	return &GeoPoint{Lat: lat, Long: long}, nil
}
