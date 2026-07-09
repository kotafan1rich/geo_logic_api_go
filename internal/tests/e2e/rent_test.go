package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validLat     = 67.6767
	validLong    = 67.6767
	validAddress = "lol"
	validInfo    = "lool"
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

	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, validLat, createdRent1.Lat)
	assert.Equal(t, validLong, createdRent1.Long)
	assert.Equal(t, validAddress, createdRent1.Address)
	assert.Equal(t, (*string)(nil), createdRent1.Info)

	resp2, err := httpAddRent(validLat, validLong, validAddress, ptr(validInfo))
	require.NoError(t, err)
	defer resp2.Body.Close()

	var createdRent2 dto.RentResponse
	err = parseBody(resp2.Body, &createdRent2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp2.StatusCode)
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
	createdRent, err := addRent(validLat, validLong, validAddress, ptr(validInfo))
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
		Info:    ptr(validInfo),
	}

	tests := []testCase{
		{
			name:             "Update lat",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: ptr(1.0)},
			expectedResponse: dto.RentResponse{Lat: 1.0, Long: validLong, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update long",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Long: ptr(1.0)},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: 1.0, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update address",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Address: ptr("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: "ll", Info: ptr(validInfo)},
		},
		{
			name:             "Update info",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Info: ptr("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: ptr("ll")},
		},
		{
			name:             "Update 2 params",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: ptr(1.0), Long: ptr(2.0)},
			expectedResponse: dto.RentResponse{Lat: 1.0, Long: 2.0, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update address with empty info",
			baseRent:         dto.CreateRentRequest{Lat: validLat, Long: validLong, Address: validAddress},
			request:          dto.UpdateRentRequest{Address: ptr("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: "ll", Info: (*string)(nil)},
		},
		{
			name:             "Update info with empty info",
			baseRent:         dto.CreateRentRequest{Lat: validLat, Long: validLong, Address: validAddress},
			request:          dto.UpdateRentRequest{Info: ptr("ll")},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: ptr("ll")},
		},
		{
			name:             "Update lat to absolute zero",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: ptr(0.0)}, // Раньше бы это поле проигнорировалось
			expectedResponse: dto.RentResponse{Lat: 0.0, Long: validLong, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update long to absolute zero",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Long: ptr(0.0)},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: 0.0, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update lat to max boundary (North Pole)",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Lat: ptr(90.0)},
			expectedResponse: dto.RentResponse{Lat: 90.0, Long: validLong, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update long to max boundary",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Long: ptr(180.0)},
			expectedResponse: dto.RentResponse{Lat: validLat, Long: 180.0, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:     "Update absolutely all fields simultaneously",
			baseRent: validRentRequest,
			request: dto.UpdateRentRequest{
				Lat:     ptr(45.5),
				Long:    ptr(120.3),
				Address: ptr("Новый проспект, дом 10"),
				Info:    ptr("Вход со двора, код 44"),
			},
			expectedResponse: dto.RentResponse{Lat: 45.5, Long: 120.3, Address: "Новый проспект, дом 10", Info: ptr("Вход со двора, код 44")},
		},
		{
			name:             "Empty patch request (no body fields sent)",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{}, // Все поля внутри равны nil
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: ptr(validInfo)},
		},
		{
			name:             "Update info from text to empty string",
			baseRent:         validRentRequest,
			request:          dto.UpdateRentRequest{Info: ptr("")}, // Затираем комментарий
			expectedResponse: dto.RentResponse{Lat: validLat, Long: validLong, Address: validAddress, Info: ptr("")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
