package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/rent/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validLat = 67.6767
const validLong = 67.6767
const validAddress = "lol"
const validInfo = "lool"

func httpAddRent(lat, long float64, address, info string) (*http.Response, error) {
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

// func addRent(lat, long float64, address, info string) (*dto.RentResponse, error) {
// 	resp, err := httpAddRent(lat, long, address, info)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var createdRent dto.RentResponse
// 	err = parseBody(resp.Body, &createdRent)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &createdRent, nil
// }

func TestE2E_RentAdd(t *testing.T) {
	clearTables(t)

	resp1, err := httpAddRent(validLat, validLong, validAddress, "")
	require.NoError(t, err)
	defer resp1.Body.Close()

	var createdRent1 dto.RentResponse
	err = parseBody(resp1.Body, &createdRent1)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, validLat, createdRent1.Lat)
	assert.Equal(t, validLong, createdRent1.Long)
	assert.Equal(t, validAddress, createdRent1.Address)
	assert.Equal(t, "", createdRent1.Info)

	resp2, err := httpAddRent(validLat, validLong, validAddress, validInfo)
	require.NoError(t, err)
	defer resp2.Body.Close()

	var createdRent2 dto.RentResponse
	err = parseBody(resp2.Body, &createdRent2)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp2.StatusCode)
	assert.Equal(t, validLat, createdRent2.Lat)
	assert.Equal(t, validLong, createdRent2.Long)
	assert.Equal(t, validAddress, createdRent2.Address)
	assert.Equal(t, validInfo, createdRent2.Info)
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
