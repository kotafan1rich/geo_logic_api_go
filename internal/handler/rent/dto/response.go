package dto

type RentResponse struct {
	ID      uint64  `json:"id"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
	Address string  `json:"address"`
	Info    string  `json:"info"`
}
