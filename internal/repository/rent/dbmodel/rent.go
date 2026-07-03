package dbmodel

import (
	"database/sql/driver"
	"fmt"

	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type DBGeoPoint model.GeoPoint

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
		return fmt.Errorf("несовместимый тип для DBGeoPoint: %T", val)
	}

	_, err := fmt.Sscanf(string(source), "POINT(%f %f)", &p.Long, &p.Lat)
	return err
}

type Rent struct {
	database.Base
	Location DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
	Address  string     `gorm:"column:address;type:varchar;size:255"`
	Info     string     `gorm:"column:info;type:varchar;size:255;default:null"`
}
