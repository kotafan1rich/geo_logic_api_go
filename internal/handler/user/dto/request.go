package dto

type CreateUserRequest struct {
	TgId uint64 `json:"tg_id"`
}

type UpdateUserRequest struct {
	Id   uint64 `json:"id"`
	TgId uint64 `json:"tg_id"`
}
