package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	clientadvertisements "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := clientadvertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

func CreateAdvertisement(t *testing.T, token string, request clientadvertisements.AdvertisementRequest, expectedStatusCode int) models.AdvertisementID {
	raw := clientadvertisements.HttpPostAdvertisement(t, token, request)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code must be %d, response: %s", expectedStatusCode, string(body))

	parsed := gjson.Get(string(body), "id").String()
	if parsed == "" {
		parsed = gjson.Get(string(body), "data.id").String()
	}
	return models.AdvertisementID{ID: parsed}
}

func GetAdvertisementByID(t *testing.T, token, adID string, expectedStatusCode int) models.Advertisement {
	raw := clientadvertisements.HttpGetAdvertisement(t, token, adID)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code must be %d, response: %s", expectedStatusCode, string(body))

	payload := gjson.Get(string(body), "data")
	if !payload.Exists() {
		payload = gjson.Parse(string(body))
	}

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal([]byte(payload.Raw), &resp), "response body: %s", string(body))
	return resp
}

func DeleteAdvertisement(t *testing.T, token, adID string, expectedStatusCode int) {
	raw := clientadvertisements.HttpDeleteAdvertisement(t, token, adID)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code must be %d, response: %s", expectedStatusCode, string(body))
}

func GetRawResponse(t *testing.T, token string, request clientadvertisements.AdvertisementRequest) (int, string) {
	raw := clientadvertisements.HttpPostAdvertisement(t, token, request)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	return raw.StatusCode, string(body)
}
