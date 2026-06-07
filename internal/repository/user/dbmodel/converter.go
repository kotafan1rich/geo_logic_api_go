package dbmodel

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToUserModel(user *model.User) *User {
	return &User{TgID: user.TgID}
}

func ToUser(userRecord *User) *model.User {
	return &model.User{
		ID:        userRecord.ID,
		TgID:      userRecord.TgID,
		CreatedAt: userRecord.CreatedAt,
		UpdatedAt: userRecord.UpdatedAt,
	}
}
