package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

func validateLat(lat float64) bool {
	if lat < minLatitude || lat > maxLatitude {
		return false
	}
	return true
}

func validateLong(long float64) bool {
	if long < minLongitude || long > maxLongitude {
		return false
	}
	return true
}

type GeoPoint struct {
	Lat  float64
	Long float64
}

func NewGeoPoint(lat, long float64) (*GeoPoint, error) {
	if !validateLat(lat) {
		return nil, errors.ErrInvalidLat
	}
	if !validateLong(long) {
		return nil, errors.ErrInvalidLong
	}
	return &GeoPoint{Lat: lat, Long: long}, nil
}

func (g *GeoPoint) UpdateLat(lat float64) error {
	if !validateLat(lat) {
		return errors.ErrInvalidLat
	}
	g.Lat = lat
	return nil
}

func (g *GeoPoint) UpdateLong(long float64) error {
	if !validateLong(long) {
		return errors.ErrInvalidLat
	}
	g.Long = long
	return nil
}
