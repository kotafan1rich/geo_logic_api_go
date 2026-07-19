package dto

type TypeResponse struct {
	ID        uint64  `json:"id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Weight    float64 `json:"weight"`
	MaxRadius uint16  `json:"max_radius"`
}
