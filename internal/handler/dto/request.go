package dto

type IDUriRequest struct {
	ID uint64 `uri:"id" binding:"required,gt=0"`
}
