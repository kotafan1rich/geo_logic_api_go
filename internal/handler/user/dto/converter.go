package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToUserResponse(user *model.User) UserResponse {
	return UserResponse{ID: user.ID, TgID: user.TgID}
}
