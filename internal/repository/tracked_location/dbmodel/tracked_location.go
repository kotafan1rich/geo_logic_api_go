package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
)

type TrackedLocation struct {
	database.Base
	UserID   uint64              `gorm:"column:user_id;type:bigint;not null;uniqueIndex"`
	Location database.DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
}
