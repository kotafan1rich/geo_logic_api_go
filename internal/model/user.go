package model

import "time"

type User struct {
	ID        uint64
	TgID      uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}
