//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/infra/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpAddInfra(lat, long float64, address string, name *string, typeID uint64) (*http.Response, error) {
	payload := dto.CreateInfraRequest{
		Lat:     lat,
		Long:    long,
		Address: address,
		Name:    name,
		TypeId:  typeID,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return httpClient.Post(infrasAPI(), "application/json", bytes.NewBuffer(jsonBytes))
}

func addInfra(lat, long float64, address string, name *string, typeID uint64) (*dto.InfraResponse, error) {
	resp, err := httpAddInfra(lat, long, address, name, typeID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createdInfra dto.InfraResponse
	err = parseBody(resp.Body, &createdInfra)
	if err != nil {
		return nil, err
	}
	return &createdInfra, nil
}

func TestE2E_InfraAdd(t *testing.T) {
	clearTables(t)

	infraType, err := addInfraType(validSlug, validName, validWeight, validMaxRadius)
	require.NoError(t, err)

	resp1, err := httpAddInfra(validLat, validLong, validAddress, nil, infraType.ID)
	require.NoError(t, err)
	defer resp1.Body.Close()

	var createdInfra1 dto.InfraResponse
	err = parseBody(resp1.Body, &createdInfra1)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, validLat, createdInfra1.Lat)
	assert.Equal(t, validLong, createdInfra1.Long)
	assert.Equal(t, validAddress, createdInfra1.Address)
	assert.Equal(t, (*string)(nil), createdInfra1.Name)
	assert.Equal(t, infraType.Name, createdInfra1.Type)

	resp2, err := httpAddInfra(validLat, validLong, validAddress, new(validName), infraType.ID)
	require.NoError(t, err)
	defer resp2.Body.Close()

	var createdInfra2 dto.InfraResponse
	err = parseBody(resp2.Body, &createdInfra2)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp2.StatusCode)
	assert.Equal(t, validLat, createdInfra2.Lat)
	assert.Equal(t, validLong, createdInfra2.Long)
	assert.Equal(t, validAddress, createdInfra2.Address)
	assert.Equal(t, validName, *createdInfra2.Name)
	assert.Equal(t, infraType.Name, createdInfra2.Type)
}

func TestE2E_InfraAdd_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача невалидного -lat (меньше -90)",
			json: `{"type_id": 1, "lat": -90.1, "long": 30.0, "address": "Тест ул."}`,
		},
		{
			name: "Передача невалидного +lat (больше 90)",
			json: `{"type_id": 1, "lat": 90.1, "long": 30.0, "address": "Тест ул."}`,
		},
		{
			name: "Отсутствие lat в запросе",
			json: `{"type_id": 1, "long": 30.0, "address": "Тест ул."}`,
		},
		{
			name: "Передача невалидного -long (меньше -180)",
			json: `{"type_id": 1, "lat": 45.0, "long": -180.1, "address": "Тест ул."}`,
		},
		{
			name: "Передача невалидного +long (больше 180)",
			json: `{"type_id": 1, "lat": 45.0, "long": 180.1, "address": "Тест ул."}`,
		},
		{
			name: "Отсутствие long в запросе",
			json: `{"type_id": 1, "lat": 45.0, "address": "Тест ул."}`,
		},
		{
			name: "Отсутствие address (передана пустая строка)",
			json: `{"type_id": 1, "lat": 45.0, "long": 30.0, "address": ""}`,
		},
		{
			name: "Отсутствие поля address в JSON",
			json: `{"type_id": 1, "lat": 45.0, "long": 30.0}`,
		},
		{
			name: "Передача пустой строки в указатель name (нарушение gt=0)",
			json: `{"type_id": 1, "lat": 45.0, "long": 30.0, "address": "Тест ул.", "name": ""}`,
		},
		{
			name: "Отсутствие type_id в запросе",
			json: `{"lat": 45.0, "long": 30.0, "address": "Тест ул."}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"type_id": 1, "lat": 45.0, "long"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Post(infrasAPI(), "application/json", bytes.NewBufferString(tc.json))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpGetInfraByID(id int64) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/%d", infrasAPI(), id))
}

func TestE2E_InfraGetByID(t *testing.T) {
	clearTables(t)

	infraType, err := addInfraType(validSlug, validName, validWeight, validMaxRadius)
	require.NoError(t, err)

	createdInfra, err := addInfra(validLat, validLong, validAddress, new(validName), infraType.ID)
	require.NoError(t, err)

	resp, err := httpGetInfraByID(int64(createdInfra.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	var parsedInfra dto.InfraResponse
	err = parseBody(resp.Body, &parsedInfra)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, validLat, createdInfra.Lat)
	assert.Equal(t, validLong, createdInfra.Long)
	assert.Equal(t, validAddress, createdInfra.Address)
	assert.Equal(t, validName, *createdInfra.Name)
	assert.Equal(t, infraType.Name, createdInfra.Type)
}

func TestE2E_InfraGetByID_NotFound(t *testing.T) {
	clearTables(t)
	resp, err := httpGetInfraByID(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_InfraGetByID_Validation(t *testing.T) {
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
			resp, err := httpClient.Get(fmt.Sprintf("%s/%s", infrasAPI(), tc.ID))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func httpUpdateInfra(id int64, payload any) (*http.Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/%d", infrasAPI(), id),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return httpClient.Do(request)
}

func TestE2E_InfraUpdate(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name             string
		baseInfra        dto.CreateInfraRequest
		request          dto.UpdateInfraRequest
		expectedResponse dto.InfraResponse
	}

	infraType, err := addInfraType(validSlug, validName, validWeight, validMaxRadius)
	require.NoError(t, err)

	newType, err := addInfraType("slug", validName, validWeight, validMaxRadius)
	require.NoError(t, err)

	validInfraRequest := dto.CreateInfraRequest{
		Lat:     validLat,
		Long:    validLong,
		Address: validAddress,
		Name:    new(validInfo),
		TypeId:  infraType.ID,
	}

	tests := []testCase{
		{
			name:      "Update lat",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Lat: new(1.0)},
			expectedResponse: dto.InfraResponse{
				Lat: 1.0, Long: validLong, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update long",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Long: new(1.0)},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: 1.0, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update address",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Address: new("Новый адрес, 22")},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: validLong, Address: "Новый адрес, 22",
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update name",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Name: new("Кафе Ромашка")},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: validLong, Address: validAddress,
				Name: new("Кафе Ромашка"), Type: infraType.Name,
			},
		},
		{
			name:      "Update type_id (Смена категории объекта)",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{TypeId: &newType.ID},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: validLong, Address: validAddress,
				Name: new(validInfo), Type: newType.Name,
			},
		},
		{
			name:      "Update lat and long simultaneously",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Lat: new(10.0), Long: new(20.0)},
			expectedResponse: dto.InfraResponse{
				Lat: 10.0, Long: 20.0, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update with empty name initially (nil name)",
			baseInfra: dto.CreateInfraRequest{Lat: validLat, Long: validLong, Address: validAddress, TypeId: infraType.ID, Name: nil},
			request:   dto.UpdateInfraRequest{Address: new("Новый Адрес")},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: validLong, Address: "Новый Адрес",
				Name: (*string)(nil), Type: infraType.Name,
			},
		},
		{
			name:      "Update lat to absolute zero",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Lat: new(0.0)},
			expectedResponse: dto.InfraResponse{
				Lat: 0.0, Long: validLong, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update long to absolute zero",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Long: new(0.0)},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: 0.0, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update lat to max boundary (North Pole)",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Lat: new(90.0)},
			expectedResponse: dto.InfraResponse{
				Lat: 90.0, Long: validLong, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update long to max boundary",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{Long: new(180.0)},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: 180.0, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Empty patch request (no body fields sent)",
			baseInfra: validInfraRequest,
			request:   dto.UpdateInfraRequest{},
			expectedResponse: dto.InfraResponse{
				Lat: validLat, Long: validLong, Address: validAddress,
				Name: new(validInfo), Type: infraType.Name,
			},
		},
		{
			name:      "Update absolutely all fields simultaneously",
			baseInfra: validInfraRequest,
			request: dto.UpdateInfraRequest{
				Lat:     new(45.5),
				Long:    new(120.3),
				Address: new("Загородный проспект, 15"),
				Name:    new("Метро Владимирская"),
				TypeId:  &newType.ID,
			},
			expectedResponse: dto.InfraResponse{
				Lat: 45.5, Long: 120.3, Address: "Загородный проспект, 15",
				Name: new("Метро Владимирская"), Type: newType.Name,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			createdInfra, err := addInfra(tc.baseInfra.Lat, tc.baseInfra.Long, tc.baseInfra.Address, tc.baseInfra.Name, tc.baseInfra.TypeId)
			require.NoError(t, err)
			tc.expectedResponse.ID = createdInfra.ID

			resp, err := httpUpdateInfra(int64(createdInfra.ID), tc.request)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedInfra dto.InfraResponse
			err = parseBody(resp.Body, &updatedInfra)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResponse, updatedInfra)
		})
	}
}

func TestE2E_InfraUpdate_NotFound(t *testing.T) {
	clearTables(t)

	infraType, err := addInfraType(validSlug, validName, validWeight, validMaxRadius)
	require.NoError(t, err)

	validInfraRequest := dto.CreateInfraRequest{
		Lat:     validLat,
		Long:    validLong,
		Address: validAddress,
		Name:    new(validInfo),
		TypeId:  infraType.ID,
	}

	resp, err := httpUpdateInfra(1, validInfraRequest)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_InfraUpdate_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача невалидного -lat (меньше -90)",
			json: `{"lat": -90.1}`,
		},
		{
			name: "Передача невалидного +lat (больше 90)",
			json: `{"lat": 90.1}`,
		},
		{
			name: "Передача невалидного +long (больше 180)",
			json: `{"long": 180.1}`,
		},
		{
			name: "Передача невалидного -long (меньше -180)",
			json: `{"long": -180.1}`,
		},
		{
			name: "Передача пустой строки в указатель name (нарушение gt=0)",
			json: `{"name": ""}`,
		},
		{
			name: "Передача пустой строки в address",
			json: `{"address": ""}`,
		},
		{
			name: "Передача нуля в type_id (невалидный ID)",
			json: `{"type_id": 0}`,
		},
		{
			name: "Сломанный синтаксис JSON (пропущена кавычка)",
			json: `{"address": "Улица`,
		},
		{
			name: "Сломанный синтаксис JSON (не закрыта фигурная скобка)",
			json: `{"lat": 45.0, "long": 34.0`,
		},
		{
			name: "Передача некорректного типа данных в координаты (строка вместо числа)",
			json: `{"lat": "сорок пять"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("%s/%d", infrasAPI(), 1),
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

func httpDeleteInfra(id int64) (*http.Response, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d", infrasAPI(), id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(request)
}

func TestE2E_InfraDelete(t *testing.T) {
	clearTables(t)

	infraType, err := addInfraType(validSlug, validName, validWeight, validMaxRadius)
	require.NoError(t, err)

	infra, err := addInfra(validLat, validLong, validAddress, new(validName), infraType.ID)
	require.NoError(t, err)

	resp1, err := httpDeleteInfra(int64(infra.ID))
	require.NoError(t, err)
	defer resp1.Body.Close()

	require.Equal(t, http.StatusNoContent, resp1.StatusCode)

	resp2, err := httpGetInfraByID(int64(infra.ID))
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}

func TestE2E_InfraDelete_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpDeleteInfra(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_InfraDelete_Validation(t *testing.T) {
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
				fmt.Sprintf("%s/%s", infrasAPI(), tc.id),
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
func httpGetInfraNear(queryParams string) (*http.Response, error) {
	url := fmt.Sprintf("%s/near?%s", infrasAPI(), queryParams)
	return httpClient.Get(url)
}

func TestE2E_InfraNear_Success(t *testing.T) {
	clearTables(t)

	// Координаты центра поиска (Центр СПБ)
	centerLat := 59.931
	centerLong := 30.354

	// 1. Категория с МАЛЕНЬКИМ радиусом (500м) -> объект в ~150м (НАЙДЕТСЯ)
	typeSmall, err := addInfraType("shop_small", "Магазин", 1.0, 500)
	require.NoError(t, err)
	nearInfra, err := addInfra(59.932, 30.355, validAddress, new(validInfo), typeSmall.ID)
	require.NoError(t, err)

	// 2. Категория со СРЕДНИМ радиусом (4000м) -> объект в ~2.5км (НАЙДЕТСЯ)
	typeMedium, err := addInfraType("cafe_mid", "Кафе", 1.0, 4000)
	require.NoError(t, err)
	midInfra, err := addInfra(59.950, 30.330, validAddress, new(validInfo), typeMedium.ID)
	require.NoError(t, err)

	// 3. Категория с ОГРАНИЧЕНИЕМ (3000м) -> объект в ~7км (НЕ НАЙДЕТСЯ, вылетел за лимит категории)
	typeFarRestricted, err := addInfraType("atm_far", "Банкомат", 1.0, 3000)
	require.NoError(t, err)
	_, err = addInfra(59.990, 30.310, validAddress, new(validInfo), typeFarRestricted.ID)
	require.NoError(t, err)

	t.Run("Поиск объектов по радиусу их категорий", func(t *testing.T) {
		queryParams := fmt.Sprintf("lat=%f&long=%f", centerLat, centerLong)
		resp, err := httpGetInfraNear(queryParams)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []dto.InfraResponse
		err = parseBody(resp.Body, &res)
		require.NoError(t, err)

		// Должно вернуться ровно 2 объекта, дальний отсекается базой по t.max_radius
		require.Len(t, res, 2)
		assert.Equal(t, nearInfra.ID, res[0].ID)
		assert.Equal(t, midInfra.ID, res[1].ID)
	})
}

func TestE2E_InfraNear_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name        string
		queryParams string
	}

	tests := []testCase{
		{
			name:        "Пропущен обязательный параметр long",
			queryParams: fmt.Sprintf("lat=%f", validLat),
		},
		{
			name:        "Пропущен обязательный параметр lat",
			queryParams: fmt.Sprintf("long=%f", validLong),
		},
		{
			name:        "Передан невалидный lat (больше 90)",
			queryParams: fmt.Sprintf("lat=90.1&long=%f", validLong),
		},
		{
			name:        "Передан невалидный lat (меньше -90)",
			queryParams: fmt.Sprintf("lat=-90.1&long=%f", validLong),
		},
		{
			name:        "Передан невалидный long (больше 180)",
			queryParams: fmt.Sprintf("lat=%f&long=180.1", validLat),
		},
		{
			name:        "Передан невалидный long (меньше -180)",
			queryParams: fmt.Sprintf("lat=%f&long=-180.1", validLat),
		},
		{
			name:        "Нечисловое значение в координатах",
			queryParams: "lat=string_here&long=30.354",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpGetInfraNear(tc.queryParams)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

