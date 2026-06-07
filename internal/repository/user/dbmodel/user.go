package dbmodel

import "github.com/kotafan1rich/geo_logic_api_go/internal/database"

type User struct {
	database.Base
	TgID uint64 `gorm:"column:tg_id;type:bigint;unique;not null"`
}
