//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func parseBody(body io.ReadCloser, dest any) error {
	err := json.NewDecoder(body).Decode(dest)
	if err != nil {
		return err
	}
	return nil
}

func addUser(tgId uint64) (*dto.UserResponse, error) {
	resp, err := httpAddUser(tgId)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createUserResponse dto.UserResponse
	err = parseBody(resp.Body, &createUserResponse)
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

func httpGetByID(id int64) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/api/user/get_by_id/%d", testServerURL, id))
}

func getByID(id int64) (*dto.UserResponse, error) {
	resp, err := httpGetByID(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user dto.UserResponse
	err = parseBody(resp.Body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func TestE2E_UserGetByID(t *testing.T) {
	clearTables(t)
	var tgId uint64 = 1
	createdUser, err := addUser(tgId)
	require.NoError(t, err)

	resp, err := httpGetByID(int64(createdUser.ID))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, createdUser.TgID, tgId)
}

func TestE2E_UserGetByID_NotFound(t *testing.T) {
	clearTables(t)
	resp, err := httpGetByID(1)
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

func httpUpdateUser(updateUser *dto.UpdateUserRequest) (*http.Response, error) {
	jsonBytes, _ := json.Marshal(updateUser)

	request, err := http.NewRequest(
		http.MethodPut,
		testServerURL+"/api/user/update",
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return httpClient.Do(request)
}

func updateUser(id uint64, tgID uint64) (*dto.UserResponse, error) {
	updateUser := &dto.UpdateUserRequest{
		ID:   id,
		TgID: tgID,
	}

	resp, err := httpUpdateUser(updateUser)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var updateUserResponse dto.UserResponse
	err = parseBody(resp.Body, &updateUserResponse)
	if err != nil {
		return nil, err
	}

	return &updateUserResponse, nil
}

func TestE2E_UserUpdate(t *testing.T) {
	clearTables(t)
	oldTgID := 1
	newTgID := 2

	newUser, err := addUser(uint64(oldTgID))
	require.NoError(t, err)

	resp, err := httpUpdateUser(&dto.UpdateUserRequest{ID: newUser.ID, TgID: uint64(newTgID)})
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedUser dto.UserResponse
	err = parseBody(resp.Body, &updatedUser)
	require.NoError(t, err)

	factUser, err := getByID(int64(newUser.ID))
	require.NoError(t, err)

	assert.Equal(t, newTgID, int(updatedUser.TgID))
	assert.Equal(t, factUser.TgID, updatedUser.TgID)
}

func TestE2E_UserUpdate_Conflict(t *testing.T) {
	clearTables(t)
	tgID1 := 1
	tgID2 := 2

	_, err := addUser(uint64(tgID1))
	require.NoError(t, err)

	newUser2, err := addUser(uint64(tgID2))
	require.NoError(t, err)

	resp, err := httpUpdateUser(&dto.UpdateUserRequest{ID: newUser2.ID, TgID: uint64(tgID1)})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestE2E_UserUpdate_Validation(t *testing.T) {
	clearTables(t)
	type testCase struct {
		name string
		json string
	}

	tests := []testCase{
		{
			"Отрицательный tgID",
			`{"id": 1, "tg_id": -1}`,
		},
		{
			"tgID не число",
			`{"id": 1, "tg_id": "lol"}`,
		},
		{
			"Отрицательный ID",
			`{"id": -1, "tg_id": 1}`,
		},
		{
			"ID не число",
			`{"id": "lol", "tg_id": 1}`,
		},
		{
			"нет ID",
			`{"tg_id": 1}`,
		},
		{
			"нет TgID",
			`{"id": 1}`,
		},
		{
			"пустой",
			"{}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jsonBytes, _ := json.Marshal(dto.UpdateUserRequest{})

			request, err := http.NewRequest(
				http.MethodPut,
				testServerURL+"/api/user/update",
				bytes.NewBuffer(jsonBytes),
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

func httpDeleteUser(id int64) (*http.Response, error) {
	request, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/user/delete/%d", testServerURL, id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(request)
}

func TestE2E_UserDelete(t *testing.T) {
	clearTables(t)

	tgID := 1
	user, err := addUser(uint64(tgID))
	require.NoError(t, err)

	resp1, err := httpDeleteUser(int64(user.ID))
	require.NoError(t, err)
	defer resp1.Body.Close()

	require.Equal(t, http.StatusNoContent, resp1.StatusCode)

	resp2, err := httpGetByID(int64(user.ID))
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)
}

func TestE2E_UserDelete_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpDeleteUser(1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_UserDelete_Validation(t *testing.T) {
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
		request, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/api/user/delete/%s", testServerURL, tc.id),
			nil,
		)
		require.NoError(t, err)

		resp, err := httpClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}
