package model

import "time"

type Event struct {
	ID      uint64
	Geopint GeoPoint
	Date    time.Time
	Info    *string
}

func NewEvent(id uint64, geopoint GeoPoint, date time.Time, info *string) *Event {
	return &Event{
		ID:      id,
		Geopint: geopoint,
		Date:    date,
		Info:    info,
	}
}
