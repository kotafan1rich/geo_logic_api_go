package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
)

type Rent struct {
	database.Base
	Location database.DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
	Address  string              `gorm:"column:address;type:varchar;size:255"`
	Info     *string             `gorm:"column:info;type:varchar;size:255;default:null"`
}
