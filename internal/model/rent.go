package model

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

func NewRent(geopoint *GeoPoint, address, info string) (*Rent, error) {
	return &Rent{
		Lat:     geopoint.Lat,
		Long:    geopoint.Long,
		Address: address,
		Info:    info,
	}, nil
}
