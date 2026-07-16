package dto

import "time"

type CreateEventRequest struct {
	Lat  float64   `json:"lat" binding:"required,gte=-90,lte=90"`
	Long float64   `json:"long" binding:"required,gte=-180,lte=180"`
	Date time.Time `json:"address" binding:"required"`
	Info *string   `json:"info,omitempty"`
}

type UpdateEventRequest struct {
	Lat  *float64   `json:"lat,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Long *float64   `json:"long,omitempty" binding:"omitempty,gte=-180,lte=180"`
	Date *time.Time `json:"date,omitempty"`
	Info *string    `json:"info,omitempty"`
}

type AvailableEventRequest struct {
	Lat    float64 `form:"lat" binding:"required,gte=-90,lte=90"`
	Long   float64 `form:"long" binding:"required,gte=-180,lte=180"`
	Radius *uint16 `form:"radius" binding:"omitempty,gt=0,lte=5000"`
}
