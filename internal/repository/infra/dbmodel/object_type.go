package dbmodel

import "github.com/kotafan1rich/geo_logic_api_go/internal/database"

type InfraType struct {
	database.Base

	Slug      string  `gorm:"column:slug;uniqueIndex;type:varchar;size:255;not null"`
	Name      string  `gorm:"column:name;type:varchar;size:255;not null"`
	Weight    float64 `gorm:"column:weight;type:float;not null"`
	MaxRadius uint16  `gorm:"column:max_radius;type:int;not null"`
}
