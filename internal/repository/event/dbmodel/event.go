package dbmodel

import (
	"time"

	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
)

type Event struct {
	database.Base

	Location database.DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
	Date     time.Time           `gorm:"column:date;type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	Info     *string             `gorm:"column:info;type:varchar;size:255"`
}
