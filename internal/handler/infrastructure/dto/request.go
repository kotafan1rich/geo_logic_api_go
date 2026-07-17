package dto

type CreateTypeRequest struct {
	Slug      string  `json:"slug" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Weight    float64 `json:"weight" binding:"gt=0,required"`
	MaxRadius uint16  `json:"max_radius" binding:"gt=0,required"`
}

type UpdateTypeRequest struct {
	Slug      *string  `json:"slug,omitempty"`
	Name      *string  `json:"name,omitempty"`
	Weight    *float64 `json:"weight,omitempty" binding:"gt=0"`
	MaxRadius *uint16  `json:"max_radius,omitempty" binding:"gt=0"`
}
