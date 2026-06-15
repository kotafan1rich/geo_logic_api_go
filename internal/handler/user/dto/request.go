package dto

type CreateUserRequest struct {
	TgID uint64 `json:"tg_id" binding:"gt=0,required"`
}

type UpdateUserRequest struct {
	TgID uint64 `json:"tg_id" binding:"gt=0,required"`
}
