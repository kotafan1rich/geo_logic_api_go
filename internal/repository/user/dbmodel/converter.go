package dbmodel

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToUserModel(user *model.User) *User {
	userModel := User{TgID: user.TgID}
	if user.ID > 0 {
		userModel.ID = user.ID
	}
	return &userModel
}

func ToUser(userRecord *User) *model.User {
	return &model.User{
		ID:   userRecord.ID,
		TgID: userRecord.TgID,
	}
}
