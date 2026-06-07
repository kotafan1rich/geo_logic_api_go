package dto

type CreateUserRequest struct {
	TgID uint64 `json:"tg_id" binding:"required"`
}

type UpdateUserRequest struct {
	ID   uint64 `json:"id" binding:"required"`
	TgID uint64 `json:"tg_id" binding:"required"`
}
