package dbmodel

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

func ToEvent(event *Event) *model.Event {
	return &model.Event{
		ID:       event.ID,
		Geopoint: model.GeoPoint(event.Location),
		Date:     event.Date,
		Info:     event.Info,
	}
}

func ToEventSlice(rents []Event) []model.Event {
	result := make([]model.Event, 0, len(rents))
	for _, event := range rents {
		result = append(result, *ToEvent(&event))
	}
	return result
}

func ToEventModel(event *model.Event) *Event {
	res := &Event{
		Location: database.DBGeoPoint(event.Geopoint),
		Date:     event.Date,
		Info:     event.Info,
	}
	res.ID = event.ID
	return res
}
