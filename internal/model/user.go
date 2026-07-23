package model

type User struct {
	ID   uint64
	TgID uint64
}

func NewUser(tgID uint64) *User {
	return &User{TgID: tgID}
}

func (u *User) UpdateUser(tgID uint64) {
	u.TgID = tgID
}
