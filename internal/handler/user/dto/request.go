package dto

type CreateUserRequest struct {
	TgID uint64 `json:"tg_id" binding:"gt=0,required"`
}

type UpdateUserRequest struct {
	ID   uint64 `json:"id" binding:"gt=0,required"`
	TgID uint64 `json:"tg_id" binding:"gt=0,required"`
}
