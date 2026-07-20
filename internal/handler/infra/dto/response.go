package dto

type TypeResponse struct {
	ID        uint64  `json:"id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Weight    float64 `json:"weight"`
	MaxRadius uint16  `json:"max_radius"`
}

type InfraResponse struct {
	ID      uint64  `json:"id"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
	Address string  `json:"address"`
	Name    *string `json:"name"`
	Type    string  `json:"type"`
}
