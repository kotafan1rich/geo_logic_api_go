package dto

type CreateTrackedLocationRequest struct {
	UserID uint64  `json:"user_id" binding:"required,gt=0"`
	Lat    float64 `json:"lat" binding:"required,gte=-90,lte=90"`
	Long   float64 `json:"long" binding:"required,gte=-180,lte=180"`
}

type UpdateTrackedLocationRequest struct {
	UserID *uint64  `json:"user_id,omitempty" binding:"omitempty,gt=0"`
	Lat    *float64 `json:"lat,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Long   *float64 `json:"long,omitempty" binding:"omitempty,gte=-180,lte=180"`
}
