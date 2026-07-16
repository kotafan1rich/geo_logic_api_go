package model

import "time"

type Event struct {
	ID      uint64
	Geopint GeoPoint
	Date    time.Time
	Info    *string
}

func NewEvent(geopoint GeoPoint, date time.Time, info *string) *Event {
	return &Event{
		Geopint: geopoint,
		Date:    date,
		Info:    info,
	}
}

func (e *Event) UpdateDate(date time.Time) {
	e.Date = date
}

func (e *Event) UpdateInfo(info *string) {
	e.Info = info
}
