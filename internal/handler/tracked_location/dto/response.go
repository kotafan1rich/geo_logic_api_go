package dto

type TrackedLocationResponse struct {
	ID     uint64  `json:"id"`
	UserID uint64  `json:"user_id"`
	Lat    float64 `json:"lat"`
	Long   float64 `json:"long"`
}
