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

	"github.com/kotafan1rich/geo_logic_api_go/internal/handler/tracked_location/dto"
)

func httpAddTrackedLocation(userID uint64, lat, long float64) (*http.Response, error) {
	payload := dto.CreateTrackedLocationRequest{UserID: userID, Lat: lat, Long: long}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return httpClient.Post(trackedLocationsAPI(), "application/json", bytes.NewBuffer(jsonBytes))
}

func addTrackedLocation(userID uint64, lat, long float64) (*dto.TrackedLocationResponse, error) {
	resp, err := httpAddTrackedLocation(userID, lat, long)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var location dto.TrackedLocationResponse
	if err = parseBody(resp.Body, &location); err != nil {
		return nil, err
	}
	return &location, nil
}

func httpGetTrackedLocation(id string) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/%s", trackedLocationsAPI(), id))
}

func httpGetTrackedLocationsByUserID(userID string) (*http.Response, error) {
	return httpClient.Get(fmt.Sprintf("%s/user/%s", trackedLocationsAPI(), userID))
}

func httpUpdateTrackedLocation(id uint64, payload any) (*http.Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/%d", trackedLocationsAPI(), id),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return httpClient.Do(req)
}

func httpDeleteTrackedLocation(id string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", trackedLocationsAPI(), id), nil)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}

func TestE2E_TrackedLocationAdd(t *testing.T) {
	clearTables(t)

	resp, err := httpAddTrackedLocation(1, validLat, validLong)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var location dto.TrackedLocationResponse
	require.NoError(t, parseBody(resp.Body, &location))
	assert.NotZero(t, location.ID)
	assert.Equal(t, uint64(1), location.UserID)
	assert.Equal(t, validLat, location.Lat)
	assert.Equal(t, validLong, location.Long)
}

func TestE2E_TrackedLocationAdd_Conflict(t *testing.T) {
	clearTables(t)

	first, err := httpAddTrackedLocation(1, validLat, validLong)
	require.NoError(t, err)
	first.Body.Close()

	second, err := httpAddTrackedLocation(1, validLat+1, validLong+1)
	require.NoError(t, err)
	defer second.Body.Close()

	assert.Equal(t, http.StatusConflict, second.StatusCode)
}

func TestE2E_TrackedLocationAdd_Validation(t *testing.T) {
	clearTables(t)

	tests := []struct {
		name string
		body string
	}{
		{name: "user_id отсутствует", body: `{"lat": 10, "long": 20}`},
		{name: "user_id равен нулю", body: `{"user_id": 0, "lat": 10, "long": 20}`},
		{name: "широта меньше минимума", body: `{"user_id": 1, "lat": -91, "long": 20}`},
		{name: "широта больше максимума", body: `{"user_id": 1, "lat": 91, "long": 20}`},
		{name: "долгота меньше минимума", body: `{"user_id": 1, "lat": 10, "long": -181}`},
		{name: "долгота больше максимума", body: `{"user_id": 1, "lat": 10, "long": 181}`},
		{name: "неверный тип", body: `{"user_id": "one", "lat": 10, "long": 20}`},
		{name: "сломанный JSON", body: `{"user_id": 1`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := httpClient.Post(trackedLocationsAPI(), "application/json", bytes.NewBufferString(tc.body))
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_TrackedLocationGetByUserID(t *testing.T) {
	clearTables(t)

	created, err := addTrackedLocation(7, validLat, validLong)
	require.NoError(t, err)

	resp, err := httpGetTrackedLocationsByUserID("7")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var locations []dto.TrackedLocationResponse
	require.NoError(t, parseBody(resp.Body, &locations))
	require.Len(t, locations, 1)
	assert.Equal(t, *created, locations[0])
}

func TestE2E_TrackedLocationGetByUserID_Empty(t *testing.T) {
	clearTables(t)

	resp, err := httpGetTrackedLocationsByUserID("1")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var locations []dto.TrackedLocationResponse
	require.NoError(t, parseBody(resp.Body, &locations))
	assert.Empty(t, locations)
}

func TestE2E_TrackedLocationGetByUserID_Validation(t *testing.T) {
	clearTables(t)

	for _, id := range []string{"0", "-1", "not-a-number"} {
		t.Run(id, func(t *testing.T) {
			resp, err := httpGetTrackedLocationsByUserID(id)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_TrackedLocationGetByID(t *testing.T) {
	clearTables(t)

	created, err := addTrackedLocation(7, validLat, validLong)
	require.NoError(t, err)

	resp, err := httpGetTrackedLocation(fmt.Sprint(created.ID))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var location dto.TrackedLocationResponse
	require.NoError(t, parseBody(resp.Body, &location))
	assert.Equal(t, *created, location)
}

func TestE2E_TrackedLocationGetByID_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpGetTrackedLocation("1")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_TrackedLocationGetByID_Validation(t *testing.T) {
	clearTables(t)

	for _, id := range []string{"0", "-1", "not-a-number"} {
		t.Run(id, func(t *testing.T) {
			resp, err := httpGetTrackedLocation(id)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_TrackedLocationUpdate(t *testing.T) {
	clearTables(t)

	created, err := addTrackedLocation(1, validLat, validLong)
	require.NoError(t, err)

	payload := dto.UpdateTrackedLocationRequest{
		UserID: new(uint64(2)),
		Lat:    new(45.5),
		Long:   new(-120.3),
	}
	resp, err := httpUpdateTrackedLocation(created.ID, payload)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updated dto.TrackedLocationResponse
	require.NoError(t, parseBody(resp.Body, &updated))
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, uint64(2), updated.UserID)
	assert.Equal(t, 45.5, updated.Lat)
	assert.Equal(t, -120.3, updated.Long)
}

func TestE2E_TrackedLocationUpdate_NotFound(t *testing.T) {
	clearTables(t)

	resp, err := httpUpdateTrackedLocation(1, dto.UpdateTrackedLocationRequest{Lat: new(10.0)})
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestE2E_TrackedLocationUpdate_Validation(t *testing.T) {
	clearTables(t)

	tests := []struct {
		name string
		id   string
		body string
	}{
		{name: "невалидный id", id: "bad", body: `{"lat": 10}`},
		{name: "нулевой user_id", id: "1", body: `{"user_id": 0}`},
		{name: "невалидная широта", id: "1", body: `{"lat": 91}`},
		{name: "невалидная долгота", id: "1", body: `{"long": -181}`},
		{name: "сломанный JSON", id: "1", body: `{"lat":`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("%s/%s", trackedLocationsAPI(), tc.id),
				bytes.NewBufferString(tc.body),
			)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestE2E_TrackedLocationDelete(t *testing.T) {
	clearTables(t)

	created, err := addTrackedLocation(1, validLat, validLong)
	require.NoError(t, err)

	resp, err := httpDeleteTrackedLocation(fmt.Sprint(created.ID))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	second, err := httpDeleteTrackedLocation(fmt.Sprint(created.ID))
	require.NoError(t, err)
	defer second.Body.Close()
	assert.Equal(t, http.StatusNotFound, second.StatusCode)
}

func TestE2E_TrackedLocationDelete_Validation(t *testing.T) {
	clearTables(t)

	resp, err := httpDeleteTrackedLocation("invalid")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
