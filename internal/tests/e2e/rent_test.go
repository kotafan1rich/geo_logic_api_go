//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent/dto"
)

func httpAddRent(lat, long float64, address string, info *string) (*http.Response, error) {
	payload := dto.CreateRentRequest{
		Lat:     lat,
		Long:    long,
		Address: address,
		Info:    info,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return httpClient.Post(rentsAPI(), "application/json", bytes.NewBuffer(jsonBytes))
}

func addRent(lat, long float64, address string, info *string) (*dto.RentResponse, error) {
	resp, err := httpAddRent(lat, long, address, info)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createdRent dto.RentResponse
	err = parseBody(resp.Body, &createdRent)
	if err != nil {
		return nil, err
	}
	return &createdRent, nil
}

func TestE2E_RentAdd(t *testing.T) {
	clearTables(t)

	resp1, err := httpAddRent(validLat, validLong, validAddress, nil)
	require.NoError(t, err)
	defer resp1.Body.Close()

	var createdRent1 dto.RentResponse
	err = parseBody(resp1.Body, &createdRent1)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, validLat, createdRent1.Lat)
	assert.Equal(t, validLong, createdRent1.Long)
	assert.Equal(t, validAddress, createdRent1.Address)
	assert.Equal(t, (*string)(nil), createdRent1.Info)

	resp2, err := httpAddRent(validLat, validLong, validAddress, new(validInfo))
	require.NoError(t, err)
	defer resp2.Body.Close()

	var createdRent2 dto.RentResponse
	err = parseBody(resp2.Body, &createdRent2)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp2.StatusCode)
	assert.Equal(t, validLat, createdRent2.Lat)
	assert.Equal(t, validLong, createdRent2.Long)
	assert.Equal(t, validAddress, createdRent2.Address)
	assert.Equal(t, validInfo, *createdRent2.Info)
}

func TestE2E_RentAdd_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача невалидного -lat",
			json: `{"lat": -91,"long": 67,"address": "lol"}`,
		},
		{
			name: "Передача невалидного +lat",
			json: `{"lat": 91,"long": 67,"address": "lol"}`,
		},
		{
			name: "Отсутствие lat",
			json: `{"long": 0,"address": "lol"}`,
		},
		{
			name: "Передача невалидного +long",
			json: `{"lat": 67,"long": 181,"address": "lol"}`,
		},
		{
			name: "Передача невалидного -long",
			json: `{"lat": 67,"long": -181,"address": "lol"}`,
		},
		{
			name: "Отсутствие long",
			json: `{"lat": 0.0, "address": "lol"}`,
		},
		{
			name: "Отсутствие address",
			json: `{"lat": 0.0,"long": 18}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"lat": 0.0,"long`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := httpClient.Post(rentsAPI(), "application/json", bytes.NewBufferString(tc.json))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpGetRentByID(id int64) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/%d", rentsAPI(), id))
}

func TestE2E_RentGetByID(t *testing.T) {
	clearTables(t)
	createdRent, err := addRent(validLat, validLong, validAddress, new(validInfo))
	require.NoError(t, err)

	resp, err := httpGetRentByID(int64(createdRent.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	var parsedRent dto.RentResponse
	err = parseBody(resp.Body, &parsedRent)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, createdRent.Lat, parsedRent.Lat)
	assert.Equal(t, createdRent.Long, parsedRent.Long)
	assert.Equal(t, createdRent.Address, parsedRent.Address)
	assert.Equal(t, createdRent.Info, parsedRent.Info)
}

func TestE2E_RentGetByID_NotFound(t *testing.T) {
	clearTables(t)
	resp, err := httpGetRentByID(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_RentGetByID_Validation(t *testing.T) {
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
			t.Parallel()
			resp, err := httpClient.Get(fmt.Sprintf("%s/%s", rentsAPI(), tc.ID))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpUpdateRent(id int64, payload any) (*http.Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/%d", rentsAPI(), id),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return httpClient.Do(request)
}

func TestE2E_RentUpdate(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name             string
		baseRent         dto.CreateRentRequest
		request          dto.UpdateRentRequest
		expectedResponse dto.RentResponse
	}

	validRentRequest := dto.CreateRentRequest{
		Lat:     validLat,
		Long:    validLong,
		Address: validAddress,
		Info:    new(validInfo),
	}

	tests := []testCase{
		{
			name:             "Update lat",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: new(1.0)},
			expectedResponse: dto.RentResponse{Lat: 1.0, Long: validLong, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update long",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Long: new(1.0)},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: 1.0, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update address",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Address: new("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: "ll", Info: new(validInfo)},
		},
		{
			name:             "Update info",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Info: new("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: new("ll")},
		},
		{
			name:             "Update 2 params",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: new(1.0), Long: new(2.0)},
			expectedResponse: dto.RentResponse{Lat: 1.0, Long: 2.0, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update address with empty info",
			baseRent:         dto.CreateRentRequest{Lat: validLat, Long: validLong, Address: validAddress},
			request:          dto.UpdateRentRequest{Address: new("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: "ll", Info: (*string)(nil)},
		},
		{
			name:             "Update info with empty info",
			baseRent:         dto.CreateRentRequest{Lat: validLat, Long: validLong, Address: validAddress},
			request:          dto.UpdateRentRequest{Info: new("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: new("ll")},
		},
		{
			name:             "Update lat to absolute zero",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: new(0.0)}, // Раньше бы это поле проигнорировалось
			expectedResponse: dto.RentResponse{Lat: 0.0, Long: validLong, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update long to absolute zero",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Long: new(0.0)},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: 0.0, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update lat to max boundary (North Pole)",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: new(90.0)},
			expectedResponse: dto.RentResponse{Lat: 90.0, Long: validLong, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update long to max boundary",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Long: new(180.0)},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: 180.0, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:     "Update absolutely all fields simultaneously",
			baseRent: validRentRequest,
			request: dto.UpdateRentRequest{
				Lat:     new(45.5),
				Long:    new(120.3),
				Address: new("Новый проспект, дом 10"),
				Info:    new("Вход со двора, код 44"),
			},
			expectedResponse: dto.RentResponse{Lat: 45.5, Long: 120.3, Address: "Новый проспект, дом 10", Info: new("Вход со двора, код 44")},
		},
		{
			name:             "Empty patch request (no body fields sent)",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{}, // Все поля внутри равны nil
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: new(validInfo)},
		},
		{
			name:             "Update info from text to empty string",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Info: new("")}, // Затираем комментарий
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: new("")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			createdRent, err := addRent(tc.baseRent.Lat, tc.baseRent.Long, tc.baseRent.Address, tc.baseRent.Info)
			require.NoError(t, err)
			tc.expectedResponse.ID = createdRent.ID

			resp, err := httpUpdateRent(int64(createdRent.ID), tc.request)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedRent dto.RentResponse
			err = parseBody(resp.Body, &updatedRent)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResponse, updatedRent)
		})
	}
}

func TestE2E_RentUpdate_NotFound(t *testing.T) {
	clearTables(t)

	validRentRequest := dto.CreateRentRequest{
		Lat:     validLat,
		Long:    validLong,
		Address: validAddress,
		Info:    new(validInfo),
	}

	resp, err := httpUpdateRent(1, validRentRequest)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_RentUpdate_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача невалидного -lat",
			json: `{"lat": -91,"long": 67,"address": "lol"}`,
		},
		{
			name: "Передача невалидного +lat",
			json: `{"lat": 91,"long": 67,"address": "lol"}`,
		},
		{
			name: "Передача невалидного +long",
			json: `{"lat": 67,"long": 181,"address": "lol"}`,
		},
		{
			name: "Передача невалидного -long",
			json: `{"lat": 67,"long": -181,"address": "lol"}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"lat": 0.0,"long`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("%s/%d", rentsAPI(), 1),
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

func httpDeleteRent(id int64) (*http.Response, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d", rentsAPI(), id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(request)
}

func TestE2E_RentDelete(t *testing.T) {
	clearTables(t)

	rent, err := addRent(validLat, validLong, validAddress, new(validInfo))
	require.NoError(t, err)

	resp1, err := httpDeleteRent(int64(rent.ID))
	require.NoError(t, err)
	defer resp1.Body.Close()

	require.Equal(t, http.StatusNoContent, resp1.StatusCode)

	resp2, err := httpGetRentByID(int64(rent.ID))
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}

func TestE2E_RentDelete_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpDeleteRent(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_RentDelete_Validation(t *testing.T) {
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
			t.Parallel()
			request, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/%s", rentsAPI(), tc.id),
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

func TestE2E_RentAvailable(t *testing.T) {
	clearTables(t) // Очищаем базу перед тестом

	// 1. Координаты центра (например, центр Питера)
	centerLat := 59.9398
	centerLong := 30.3146

	// 2. Создаем тестовые точки в базе данных через вашу функцию addRent
	// Точка А: Близко к центру (в районе 3.5 км)
	rentClose, err := addRent(59.9311, 30.3609, "Рядом с центром", new("Доступно"))
	require.NoError(t, err)

	// Точка Б: Очень далеко (Выборг, ~120 км)
	_, err = addRent(60.7102, 28.7469, "Очень далеко", new("Доступно"))
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
			expectID:      rentClose.ID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			url := fmt.Sprintf("%s/available?lat=%f&long=%f&radius=%d",
				rentsAPI(), centerLat, centerLong, tc.searchRadius,
			)

			resp, err := httpClient.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var foundRents []dto.RentResponse
			err = parseBody(resp.Body, &foundRents)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedCount, len(foundRents), "Количество найденных объектов не совпадает")

			if tc.expectedCount == 1 {
				assert.Equal(t, tc.expectID, foundRents[0].ID, "База вернула не тот объект")
			}
		})
	}
}

func TestE2E_RentAvailable_Validation(t *testing.T) {
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
			t.Parallel()
			url := fmt.Sprintf("%s/available?%s", rentsAPI(), tc.queryParams)

			resp, err := httpClient.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"Для кейса '%s' ожидался статус 400, но сервер вернул %d", tc.name, resp.StatusCode)
		})
	}
}
