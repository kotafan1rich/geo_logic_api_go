package model

import "time"

type Event struct {
	ID       uint64
	Geopoint GeoPoint
	Date     time.Time
	Info     *string
}

func NewEvent(geopoint GeoPoint, date time.Time, info *string) *Event {
	return &Event{
		Geopoint: geopoint,
		Date:     date,
		Info:     info,
	}
}

func (e *Event) UpdateDate(date time.Time) {
	e.Date = date
}

func (e *Event) UpdateInfo(info *string) {
	e.Info = info
}

func (e *Event) Update(lat, long *float64, date *time.Time, info *string) error {
	if lat != nil {
		err := e.Geopoint.UpdateLat(*lat)
		if err != nil {
			return err
		}
	}

	if long != nil {
		err := e.Geopoint.UpdateLong(*long)
		if err != nil {
			return err
		}
	}

	if date != nil {
		e.UpdateDate(*date)
	}

	if info != nil {
		e.UpdateInfo(info)
	}

	return nil
}
