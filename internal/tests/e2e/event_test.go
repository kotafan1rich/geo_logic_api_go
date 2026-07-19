//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/event/dto"
)

var validDate = time.Date(2026, time.July, 16, 12, 0, 0, 0, time.UTC)

func httpAddEvent(lat, long float64, date time.Time, info *string) (*http.Response, error) {
	payload := dto.CreateEventRequest{
		Lat:  lat,
		Long: long,
		Date: date,
		Info: info,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return httpClient.Post(eventsAPI(), "application/json", bytes.NewBuffer(jsonBytes))
}

func addEvent(lat, long float64, date time.Time, info *string) (*dto.EventResponse, error) {
	resp, err := httpAddEvent(lat, long, date, info)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createdEvent dto.EventResponse
	err = parseBody(resp.Body, &createdEvent)
	if err != nil {
		return nil, err
	}
	return &createdEvent, nil
}

func TestE2E_EventAdd(t *testing.T) {
	clearTables(t)

	resp1, err := httpAddEvent(validLat, validLong, validDate, nil)
	require.NoError(t, err)
	defer resp1.Body.Close()

	var createdEvent1 dto.EventResponse
	err = parseBody(resp1.Body, &createdEvent1)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, validLat, createdEvent1.Lat)
	assert.Equal(t, validLong, createdEvent1.Long)
	assert.Equal(t, validDate, createdEvent1.Date)
	assert.Equal(t, (*string)(nil), createdEvent1.Info)

	resp2, err := httpAddEvent(validLat, validLong, validDate, new(validInfo))
	require.NoError(t, err)
	defer resp2.Body.Close()

	var createdEvent2 dto.EventResponse
	err = parseBody(resp2.Body, &createdEvent2)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp2.StatusCode)
	assert.Equal(t, validLat, createdEvent2.Lat)
	assert.Equal(t, validLong, createdEvent2.Long)
	assert.Equal(t, validDate, createdEvent2.Date)
	assert.Equal(t, validInfo, *createdEvent2.Info)
}

func TestE2E_EventAdd_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача невалидного -lat",
			json: `{"lat": -91,"long": 67,"date": "lol"}`,
		},
		{
			name: "Передача невалидного +lat",
			json: `{"lat": 91,"long": 67,"date": "lol"}`,
		},
		{
			name: "Отсутствие lat",
			json: `{"long": 0,"date": "lol"}`,
		},
		{
			name: "Передача невалидного +long",
			json: `{"lat": 67,"long": 181,"date": "lol"}`,
		},
		{
			name: "Передача невалидного -long",
			json: `{"lat": 67,"long": -181,"date": "lol"}`,
		},
		{
			name: "Отсутствие long",
			json: `{"lat": 0.0, "date": "lol"}`,
		},
		{
			name: "Отсутствие date",
			json: `{"lat": 0.0,"long": 18}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"lat": 0.0,"long`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Post(eventsAPI(), "application/json", bytes.NewBufferString(tc.json))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpGetEventByID(id int64) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/%d", eventsAPI(), id))
}

func TestE2E_EventGetByID(t *testing.T) {
	clearTables(t)
	createdEvent, err := addEvent(validLat, validLong, validDate, new(validInfo))
	require.NoError(t, err)

	resp, err := httpGetEventByID(int64(createdEvent.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	var parsedEvent dto.EventResponse
	err = parseBody(resp.Body, &parsedEvent)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, createdEvent.Lat, parsedEvent.Lat)
	assert.Equal(t, createdEvent.Long, parsedEvent.Long)
	assert.Equal(t, createdEvent.Date, parsedEvent.Date)
	assert.Equal(t, createdEvent.Info, parsedEvent.Info)
}

func TestE2E_EventGetByID_NotFound(t *testing.T) {
	clearTables(t)
	resp, err := httpGetEventByID(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_EventGetByID_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		ID   string
	}

	tests := []testCase{
		{
			"Отрицательный ID",
			"-1",
		},
		{
			"Не число",
			"one",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Get(fmt.Sprintf("%s/%s", eventsAPI(), tc.ID))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpUpdateEvent(id int64, payload any) (*http.Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/%d", eventsAPI(), id),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return httpClient.Do(request)
}

func TestE2E_EventUpdate(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name             string
		baseEvent        dto.CreateEventRequest
		request          dto.UpdateEventRequest
		expectedResponse dto.EventResponse
	}

	validEventRequest := dto.CreateEventRequest{
		Lat:  validLat,
		Long: validLong,
		Date: validDate,
		Info: new(validInfo),
	}

	tests := []testCase{
		{
			name:             "Update lat",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Lat: new(1.0)},
			expectedResponse: dto.EventResponse{Lat: 1.0, Long: validLong, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update long",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Long: new(1.0)},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: 1.0, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update date",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Date: new(time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC))},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: validLong, Date: time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC), Info: new(validInfo)},
		},
		{
			name:             "Update info",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Info: new("ll")},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: validLong, Date: validDate, Info: new("ll")},
		},
		{
			name:             "Update 2 params",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Lat: new(1.0), Long: new(2.0)},
			expectedResponse: dto.EventResponse{Lat: 1.0, Long: 2.0, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update date with empty info",
			baseEvent:        dto.CreateEventRequest{Lat: validLat, Long: validLong, Date: validDate},
			request:          dto.UpdateEventRequest{Date: new(time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC))},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: validLong, Date: time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC), Info: (*string)(nil)},
		},
		{
			name:             "Update info with empty info",
			baseEvent:        dto.CreateEventRequest{Lat: validLat, Long: validLong, Date: validDate},
			request:          dto.UpdateEventRequest{Info: new("ll")},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: validLong, Date: validDate, Info: new("ll")},
		},
		{
			name:             "Update lat to absolute zero",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Lat: new(0.0)}, // Раньше бы это поле проигнорировалось
			expectedResponse: dto.EventResponse{Lat: 0.0, Long: validLong, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update long to absolute zero",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Long: new(0.0)},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: 0.0, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update lat to max boundary (North Pole)",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Lat: new(90.0)},
			expectedResponse: dto.EventResponse{Lat: 90.0, Long: validLong, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update long to max boundary",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Long: new(180.0)},
			expectedResponse: dto.EventResponse{Lat: validLat, Long: 180.0, Date: validDate, Info: new(validInfo)},
		},
		{
			name:      "Update absolutely all fields simultaneously",
			baseEvent: validEventRequest,
			request: dto.UpdateEventRequest{
				Lat:  new(45.5),
				Long: new(120.3),
				Date: new(time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC)),
				Info: new("Вход со двора, код 44"),
			},
			expectedResponse: dto.EventResponse{Lat: 45.5, Long: 120.3, Date: time.Date(2026, time.July, 17, 12, 0, 0, 0, time.UTC), Info: new("Вход со двора, код 44")},
		},
		{
			name:             "Empty patch request (no body fields sent)",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{}, // Все поля внутри равны nil
			expectedResponse: dto.EventResponse{Lat: validLat, Long: validLong, Date: validDate, Info: new(validInfo)},
		},
		{
			name:             "Update info from text to empty string",
			baseEvent:        validEventRequest,
			request:          dto.UpdateEventRequest{Info: new("")}, // Затираем комментарий
			expectedResponse: dto.EventResponse{Lat: validLat, Long: validLong, Date: validDate, Info: new("")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			createdEvent, err := addEvent(tc.baseEvent.Lat, tc.baseEvent.Long, tc.baseEvent.Date, tc.baseEvent.Info)
			require.NoError(t, err)
			tc.expectedResponse.ID = createdEvent.ID

			resp, err := httpUpdateEvent(int64(createdEvent.ID), tc.request)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedEvent dto.EventResponse
			err = parseBody(resp.Body, &updatedEvent)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResponse, updatedEvent)
		})
	}
}

func TestE2E_EventUpdate_NotFound(t *testing.T) {
	clearTables(t)

	validEventRequest := dto.CreateEventRequest{
		Lat:  validLat,
		Long: validLong,
		Date: validDate,
		Info: new(validInfo),
	}

	resp, err := httpUpdateEvent(1, validEventRequest)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_EventUpdate_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача невалидного -lat",
			json: `{"lat": -91,"long": 67,"date": "lol"}`,
		},
		{
			name: "Передача невалидного +lat",
			json: `{"lat": 91,"long": 67,"date": "lol"}`,
		},
		{
			name: "Передача невалидного +long",
			json: `{"lat": 67,"long": 181,"date": "lol"}`,
		},
		{
			name: "Передача невалидного -long",
			json: `{"lat": 67,"long": -181,"date": "lol"}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"lat": 0.0,"long`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("%s/%d", eventsAPI(), 1),
				bytes.NewBufferString(tc.json),
			)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			resp, err := httpClient.Do(request)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpDeleteEvent(id int64) (*http.Response, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d", eventsAPI(), id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(request)
}

func TestE2E_EventDelete(t *testing.T) {
	clearTables(t)

	event, err := addEvent(validLat, validLong, validDate, new(validInfo))
	require.NoError(t, err)

	resp1, err := httpDeleteEvent(int64(event.ID))
	require.NoError(t, err)
	defer resp1.Body.Close()

	require.Equal(t, http.StatusNoContent, resp1.StatusCode)

	resp2, err := httpGetEventByID(int64(event.ID))
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}

func TestE2E_EventDelete_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpDeleteEvent(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_EventDelete_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		id   string
	}

	tests := []testCase{
		{
			name: "Передача отрицательного",
			id:   "-1",
		},
		{
			name: "Передача строки вместо числа",
			id:   "lol",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/%s", eventsAPI(), tc.id),
				nil,
			)
			require.NoError(t, err)

			resp, err := httpClient.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_EventNear(t *testing.T) {
	clearTables(t) // Очищаем базу перед тестом

	// 1. Координаты центра (например, центр Питера)
	centerLat := 59.9398
	centerLong := 30.3146

	// 2. Создаем тестовые точки в базе данных через вашу функцию addEvent
	// Точка А: Близко к центру (в районе 3.5 км)
	eventClose, err := addEvent(59.9311, 30.3609, validDate, new("Доступно"))
	require.NoError(t, err)

	// Точка Б: Очень далеко (Выборг, ~120 км)
	_, err = addEvent(60.7102, 28.7469, validDate, new("Доступно"))
	require.NoError(t, err)

	type testCase struct {
		name          string
		searchRadius  uint64
		expectedCount int
		expectID      uint64
	}

	tests := []testCase{
		{
			name:          "Маленький радиус (1 км) - ничего не должно найти",
			searchRadius:  1000,
			expectedCount: 0,
		},
		{
			name:          "Средний радиус (5 км) - должен найти только близкую точку",
			searchRadius:  5000,
			expectedCount: 1,
			expectID:      eventClose.ID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/near?lat=%f&long=%f&radius=%d",
				eventsAPI(), centerLat, centerLong, tc.searchRadius,
			)

			resp, err := httpClient.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var foundEvents []dto.EventResponse
			err = parseBody(resp.Body, &foundEvents)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedCount, len(foundEvents), "Количество найденных объектов не совпадает")

			if tc.expectedCount == 1 {
				assert.Equal(t, tc.expectID, foundEvents[0].ID, "База вернула не тот объект")
			}
		})
	}
}

func TestE2E_EventNear_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name        string
		queryParams string
	}

	tests := []testCase{
		{
			name:        "Пропущена широта (lat)",
			queryParams: "long=67.067&radius=5000",
		},
		{
			name:        "Пропущена долгота (long)",
			queryParams: "lat=67.067&radius=5000",
		},
		{
			name:        "Невалидная широта (слишком большая)",
			queryParams: "lat=91.0&long=67.067&radius=5000",
		},
		{
			name:        "Невалидная широта (слишком маленькая)",
			queryParams: "lat=-91.0&long=67.067&radius=5000",
		},
		{
			name:        "Невалидная долгота (слишком большая)",
			queryParams: "lat=67.067&long=181.0&radius=5000",
		},
		{
			name:        "Невалидная долгота (слишком маленькая)",
			queryParams: "lat=67.067&long=-181.0&radius=5000",
		},
		{
			name:        "Вместо чисел переданы буквы в lat",
			queryParams: "lat=not_a_number&long=67.067&radius=5000",
		},
		{
			name:        "Вместо чисел переданы буквы в long",
			queryParams: "lat=67.067&long=text&radius=5000",
		},
		{
			name:        "Радиус равен нулю (нарушение бизнес-логики сервиса)",
			queryParams: "lat=67.067&long=67.067&radius=0",
		},
		{
			name:        "Радиус больше максимального (например, maxRadius = 65000, передаем 70000)",
			queryParams: "lat=67.067&long=67.067&radius=70000",
		},
		{
			name:        "Передан отрицательный радиус (uint16 сломается при парсинге)",
			queryParams: "lat=67.067&long=67.067&radius=-100",
		},
		{
			name:        "Вместо числа в радиус переданы буквы",
			queryParams: "lat=67.067&long=67.067&radius=kilometers",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/near?%s", eventsAPI(), tc.queryParams)

			resp, err := httpClient.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"Для кейса '%s' ожидался статус 400, но сервер вернул %d", tc.name, resp.StatusCode)
		})
	}
}
