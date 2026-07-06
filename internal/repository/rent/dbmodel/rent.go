package dbmodel

import (
	"database/sql/driver"
	"fmt"

	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type DBGeoPoint struct {
	Lat  float64
	Long float64
}

func (p *DBGeoPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Long, p.Lat), nil
}

func (p *DBGeoPoint) Scan(val any) error {
	var source []byte
	switch v := val.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return fmt.Errorf("conflict type for DBGeoPoint: %T", val)
	}

	geomt, err := wkb.Unmarshal(source)
	if err != nil {
		return fmt.Errorf("error decode WKB: %w", err)
	}

	point := geomt.FlatCoords()

	if len(point) >= 2 {
		p.Long = point[0]
		p.Lat = point[1]
		return nil
	}
	return fmt.Errorf("Failed to retrieve coordinates from the geometry")
}

type Rent struct {
	database.Base
	Location DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
	Address  string     `gorm:"column:address;type:varchar;size:255"`
	Info     string     `gorm:"column:info;type:varchar;size:255;default:null"`
}
