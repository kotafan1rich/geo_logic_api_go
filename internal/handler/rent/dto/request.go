package dto

type CreateRentRequest struct {
	Lat     float64 `json:"lat" binding:"required,gte=-90,lte=90"`
	Long    float64 `json:"long" binding:"required,gte=-180,lte=180"`
	Address string  `json:"address" binding:"required"`
	Info    string  `json:"info,omitempty"`
}

type UpdateRentRequest struct {
	Lat     float64 `json:"lat,omitempty" binding:"gte=-90,lte=90"`
	Long    float64 `json:"long,omitempty" binding:"gte=-180,lte=180"`
	Address string  `json:"address,omitempty"`
	Info    string  `json:"info,omitempty"`
}
