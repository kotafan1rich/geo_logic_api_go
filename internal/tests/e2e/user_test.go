//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/user/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addUser(tg_id uint64) (*http.Response, error) {
	payload := dto.CreateUserRequest{TgId: tg_id}

	jsonBytes, _ := json.Marshal(payload)

	return httpClient.Post(testServerURL+"/api/user/create", "application/json", bytes.NewBuffer(jsonBytes))
}

func TestE2E_UserAdd(t *testing.T) {
	clearTables(t)
	resp, err := addUser(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestE2E_UserAdd_Conflict(t *testing.T) {
	clearTables(t)
	var tgId uint64 = 1

	resp1, err := addUser(tgId)
	require.NoError(t, err)
	defer resp1.Body.Close()

	resp2, err := addUser(tgId)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusConflict, resp2.StatusCode)
}

func TestE2E_UserAdd_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			name: "Передача отрицательного",
			json: `{"tg_id": -1}`,
		},
		{
			name: "Пропуск обязательного поля (id отсутствует)",
			json: `{}`,
		},
		{
			name: "Передача строки вместо числа",
			json: `{"tg_id": "string_instead_of_number"}`,
		},
		{
			name: "Сломанный синтаксис JSON",
			json: `{"tg_id": `,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Post(
				testServerURL+"/api/user/create",
				"application/json",
				bytes.NewBufferString(tc.json),
			)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
