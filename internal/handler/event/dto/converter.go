package dto

import "github.com/kotafan1rich/geo_logic_api_go/internal/model"

func ToEventResponse(event *model.Event) *EventResponse {
	return &EventResponse{
		ID:   event.ID,
		Lat:  event.Geopoint.Lat,
		Long: event.Geopoint.Long,
		Date: event.Date,
		Info: event.Info,
	}
}

func ToEventResponseSlice(events []model.Event) []*EventResponse {
	result := make([]*EventResponse, 0, len(events))
	for _, event := range events {
		result = append(result, ToEventResponse(&event))
	}
	return result
}
