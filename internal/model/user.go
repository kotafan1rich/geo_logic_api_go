package model


type User struct {
	Base
	TgId uint64 `gorm:"column:tg_id;type:bigint;unique;not null"`
}
