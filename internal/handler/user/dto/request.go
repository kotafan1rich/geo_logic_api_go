package dto

type CreateUserRequest struct {
	TgId uint64 `json:"tg_id" binding:"required"`
}

type UpdateUserRequest struct {
	Id   uint64 `json:"id" binding:"required"`
	TgId uint64 `json:"tg_id" binding:"required"`
}
