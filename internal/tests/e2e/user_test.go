//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/user/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpAddUser(tgId uint64) (*http.Response, error) {
	payload := dto.CreateUserRequest{TgID: tgId}

	jsonBytes, _ := json.Marshal(payload)

	return httpClient.Post(testServerURL+"/api/user/create", "application/json", bytes.NewBuffer(jsonBytes))
}

func addUser(tgId uint64) (*dto.UserResponse, error) {
	resp, err := httpAddUser(tgId)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createUserResponse dto.UserResponse

	err = json.NewDecoder(resp.Body).Decode(&createUserResponse)
	if err != nil {
		return nil, err
	}

	return &createUserResponse, nil
}

func TestE2E_UserAdd(t *testing.T) {
	clearTables(t)
	resp, err := httpAddUser(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestE2E_UserAdd_Conflict(t *testing.T) {
	clearTables(t)
	var tgId uint64 = 1

	resp1, err := httpAddUser(tgId)
	require.NoError(t, err)
	defer resp1.Body.Close()

	resp2, err := httpAddUser(tgId)
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

func getByID(id int64) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/api/user/get_by_id/%d", testServerURL, id))
}

func TestE2E_UserGetByID(t *testing.T) {
	clearTables(t)
	var tgId uint64 = 1
	createdUser, err := addUser(tgId)
	require.NoError(t, err)

	resp, err := getByID(int64(createdUser.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, createdUser.TgID, tgId)
}

func TestE2E_UserGetByID_NotFound(t *testing.T) {
	clearTables(t)
	resp, err := getByID(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_UserGetByID_Validation(t *testing.T) {
	clearTables(t)

	type testCase struct {
		name string
		tgId string
	}

	tests := []testCase{
		{
			"Отрицательный tgId",
			"-1",
		},
		{
			"Не число",
			"one",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Get(fmt.Sprintf("%s/api/user/get_by_id/%s", testServerURL, tc.tgId))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
