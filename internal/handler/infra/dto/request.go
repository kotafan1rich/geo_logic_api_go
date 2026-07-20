package dto

type CreateTypeRequest struct {
	Slug      string  `json:"slug" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Weight    float64 `json:"weight" binding:"gt=0,required"`
	MaxRadius uint16  `json:"max_radius" binding:"gt=0,required"`
}

type UpdateTypeRequest struct {
	Slug      *string  `json:"slug,omitempty" binding:"omitempty,gt=0"`
	Name      *string  `json:"name,omitempty" binding:"omitempty,gt=0"`
	Weight    *float64 `json:"weight,omitempty" binding:"omitempty,gt=0"`
	MaxRadius *uint16  `json:"max_radius,omitempty" binding:"omitempty,gt=0"`
}

type CreateInfraRequest struct {
	Lat     float64 `json:"lat" binding:"required,gte=-90,lte=90"`
	Long    float64 `json:"long" binding:"required,gte=-180,lte=180"`
	Address string  `json:"address" binding:"required"`
	Name    *string `json:"name,omitempty" binding:"omitempty,gt=0"`
	TypeId  uint64  `json:"type_id" binding:"required"`
}

type UpdateInfraRequest struct {
	Lat     *float64 `json:"lat,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Long    *float64 `json:"long,omitempty" binding:"omitempty,gte=-180,lte=180"`
	Address *string  `json:"address,omitempty" binding:"omitempty,gt=0"`
	Name    *string  `json:"name,omitempty" binding:"omitempty,gt=0"`
	TypeId  *uint64  `json:"type_id,omitempty" binding:"omitempty,gt=0"`
}

type NearEventRequest struct {
	Lat  float64 `form:"lat" binding:"required,gte=-90,lte=90"`
	Long float64 `form:"long" binding:"required,gte=-180,lte=180"`
}
