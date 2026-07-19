//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/infrastructure/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpAddInfraType(slug, name string, weight float64, maxRadius uint16) (*http.Response, error) {
	payload := dto.CreateTypeRequest{
		Slug:      slug,
		Name:      name,
		Weight:    weight,
		MaxRadius: maxRadius,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return httpClient.Post(infraTypesAPI(), "application/json", bytes.NewBuffer(jsonBytes))
}

// Вспомогательная функция для добавления типа и возврата распарсенного DTO
func addInfraType(slug, name string, weight float64, maxRadius uint16) (*dto.TypeResponse, error) {
	resp, err := httpAddInfraType(slug, name, weight, maxRadius)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createdType dto.TypeResponse
	err = parseBody(resp.Body, &createdType)
	if err != nil {
		return nil, err
	}
	return &createdType, nil
}

// Вспомогательная функция для HTTP GET запроса по ID
func httpGetInfraTypeByID(id uint64) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/%d", infraTypesAPI(), id))
}

func httpUpdateInfraType(id uint64, payload any) (*http.Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/%d", infraTypesAPI(), id),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return httpClient.Do(request)
}

func httpDeleteInfraType(id uint64) (*http.Response, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/%d", infraTypesAPI(), id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(request)
}

// === ТЕСТЫ ===

// 1. Позитивный тест создания типа инфраструктуры
func TestE2E_InfraTypeAdd(t *testing.T) {
	clearTables(t)

	resp, err := httpAddInfraType("subway", "Метро", 1.5, 500)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdType dto.TypeResponse
	err = parseBody(resp.Body, &createdType)
	require.NoError(t, err)

	assert.NotZero(t, createdType.ID)
	assert.Equal(t, "subway", createdType.Slug)
	assert.Equal(t, "Метро", createdType.Name)
	assert.Equal(t, 1.5, createdType.Weight)
	assert.Equal(t, uint16(500), createdType.MaxRadius)
}

// 2. Тест валидации входящего JSON при создании
func TestE2E_InfraTypeAdd_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Отсутствие slug",
			json: `{"name": "Метро", "weight": 1.5, "max_radius": 500}`,
		},
		{
			name: "Отсутствие name",
			json: `{"slug": "subway", "weight": 1.5, "max_radius": 500}`,
		},
		{
			name: "Отрицательный weight",
			json: `{"slug": "subway", "name": "Метро", "weight": -1.2, "max_radius": 500}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"slug": "subway", "name"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Post(infraTypesAPI(), "application/json", bytes.NewBufferString(tc.json))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

// 3. Позитивный тест получения типа по ID
func TestE2E_InfraTypeGetByID(t *testing.T) {
	clearTables(t)

	// Создаем тестовую запись
	createdType, err := addInfraType("cafe", "Кафе", 1.0, 300)
	require.NoError(t, err)

	// Делаем GET запрос по созданному ID
	resp, err := httpGetInfraTypeByID(createdType.ID)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsedType dto.TypeResponse
	err = parseBody(resp.Body, &parsedType)
	require.NoError(t, err)

	assert.Equal(t, createdType.ID, parsedType.ID)
	assert.Equal(t, createdType.Slug, parsedType.Slug)
	assert.Equal(t, createdType.Name, parsedType.Name)
	assert.Equal(t, createdType.Weight, parsedType.Weight)
	assert.Equal(t, createdType.MaxRadius, parsedType.MaxRadius)
}

// 4. Тест получения несуществующего типа
func TestE2E_InfraTypeGetByID_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpGetInfraTypeByID(9999) // Заведомо отсутствующий ID
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// 5. Тест валидации ID в URL строке
func TestE2E_InfraTypeGetByID_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		ID   string
	}

	tests := []testCase{
		{
			"Отрицательный ID",
			"-5",
		},
		{
			"Не числовой ID",
			"bad-id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Get(fmt.Sprintf("%s/%s", infraTypesAPI(), tc.ID))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_InfraTypeUpdate(t *testing.T) {
	type testCase struct {
		name             string
		baseType         dto.CreateTypeRequest
		request          dto.UpdateTypeRequest
		expectedResponse dto.TypeResponse
	}

	validInfraTypeRequest := dto.CreateTypeRequest{
		Slug:      "subway",
		Name:      "Метро",
		Weight:    1.5,
		MaxRadius: 500,
	}

	tests := []testCase{
		{
			name:             "Update slug",
			baseType:         validInfraTypeRequest,
			request:          dto.UpdateTypeRequest{Slug: new("tram")},
			expectedResponse: dto.TypeResponse{Slug: "tram", Name: "Метро", Weight: 1.5, MaxRadius: 500},
		},
		{
			name:             "Update name",
			baseType:         validInfraTypeRequest,
			request:          dto.UpdateTypeRequest{Name: new("Трамвай")},
			expectedResponse: dto.TypeResponse{Slug: "subway", Name: "Трамвай", Weight: 1.5, MaxRadius: 500},
		},
		{
			name:             "Update weight",
			baseType:         validInfraTypeRequest,
			request:          dto.UpdateTypeRequest{Weight: new(float64(2.2))},
			expectedResponse: dto.TypeResponse{Slug: "subway", Name: "Метро", Weight: 2.2, MaxRadius: 500},
		},
		{
			name:             "Update max radius",
			baseType:         validInfraTypeRequest,
			request:          dto.UpdateTypeRequest{MaxRadius: new(uint16(700))},
			expectedResponse: dto.TypeResponse{Slug: "subway", Name: "Метро", Weight: 1.5, MaxRadius: 700},
		},
		{
			name:             "Update 2 params",
			baseType:         validInfraTypeRequest,
			request:          dto.UpdateTypeRequest{Slug: new("cafe"), Name: new("Кафе")},
			expectedResponse: dto.TypeResponse{Slug: "cafe", Name: "Кафе", Weight: 1.5, MaxRadius: 500},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clearTables(t)
			createdType, err := addInfraType(tc.baseType.Slug, tc.baseType.Name, tc.baseType.Weight, tc.baseType.MaxRadius)
			require.NoError(t, err)
			tc.expectedResponse.ID = createdType.ID

			resp, err := httpUpdateInfraType(createdType.ID, tc.request)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)

			var updatedType dto.TypeResponse
			err = parseBody(resp.Body, &updatedType)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResponse, updatedType)
		})
	}
}

func TestE2E_InfraTypeUpdate_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpUpdateInfraType(1, dto.UpdateTypeRequest{Slug: new("tram")})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_InfraTypeUpdate_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		id   string
		json string
	}

	tests := []testCase{
		{
			name: "Отрицательный weight",
			id:   "1",
			json: `{"weight": -1}`,
		},
		{
			name: "Пустой slug",
			id:   "1",
			json: `{"slug": ""}`,
		},
		{
			name: "Пустой name",
			id:   "1",
			json: `{"name": ""}`,
		},
		{
			name: "Отрицательный ID",
			id:   "-1",
			json: `{"name": "Кафе"}`,
		},
		{
			name: "ID не число",
			id:   "abc",
			json: `{"name": "Кафе"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("%s/%s", infraTypesAPI(), tc.id),
				bytes.NewBufferString(tc.json),
			)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			resp, err := httpClient.Do(request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_InfraTypeDelete(t *testing.T) {
	clearTables(t)

	createdType, err := addInfraType("subway", "Метро", 1.5, 500)
	require.NoError(t, err)

	resp1, err := httpDeleteInfraType(createdType.ID)
	require.NoError(t, err)
	defer resp1.Body.Close()

	require.Equal(t, http.StatusNoContent, resp1.StatusCode)

	resp2, err := httpGetInfraTypeByID(createdType.ID)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}

func TestE2E_InfraTypeDelete_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpDeleteInfraType(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_InfraTypeDelete_Validation(t *testing.T) {
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
				fmt.Sprintf("%s/%s", infraTypesAPI(), tc.id),
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
