package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
)

type InfrastructureObject struct {
	database.Base

	Location database.DBGeoPoint `gorm:"column:location;type:geometry(Point,4326);not null"`
	Address  string              `gorm:"column:address;type:varchar;size:255"`
	Name     *string             `gorm:"column:name;type:varchar;size:255;default:null"`

	TypeID uint64             `gorm:"column:type_id;type:bigint;not null;index"`
	Type   InfrastructureType `gorm:"column:type;foreign_key:TypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
