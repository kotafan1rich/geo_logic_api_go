package dto

import "time"

type EventResponse struct {
	ID   uint64    `json:"id"`
	Lat  float64   `json:"lat"`
	Long float64   `json:"long"`
	Date time.Time `json:"date"`
	Info *string   `json:"info"`
}
