package favorites

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	clientfavorites "api-tests-template/internal/client/http/favorites"
)

func readWithRetry(t *testing.T, call func() *http.Response) (int, string) {
	for attempt := 0; ; attempt++ {
		raw := call()
		body, err := io.ReadAll(raw.Body)
		_ = raw.Body.Close()
		require.NoError(t, err)

		if raw.StatusCode != http.StatusServiceUnavailable || attempt == 2 {
			return raw.StatusCode, string(body)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func AddFavorite(t *testing.T, token, adID string, expectedStatusCode int) string {
	status, body := readWithRetry(t, func() *http.Response {
		return clientfavorites.HttpPostFavorite(t, token, adID)
	})

	require.Equalf(t, expectedStatusCode, status,
		"HTTP status code must be %d, response: %s", expectedStatusCode, body)
	return body
}

func AddFavoriteRaw(t *testing.T, token, adID string) (int, string) {
	return readWithRetry(t, func() *http.Response {
		return clientfavorites.HttpPostFavorite(t, token, adID)
	})
}

func DeleteFavorite(t *testing.T, token, adID string, expectedStatusCode int) string {
	status, body := readWithRetry(t, func() *http.Response {
		return clientfavorites.HttpDeleteFavorite(t, token, adID)
	})

	require.Equalf(t, expectedStatusCode, status,
		"HTTP status code must be %d, response: %s", expectedStatusCode, body)
	return body
}

func DeleteFavoriteRaw(t *testing.T, token, adID string) (int, string) {
	return readWithRetry(t, func() *http.Response {
		return clientfavorites.HttpDeleteFavorite(t, token, adID)
	})
}

func GetFavorites(t *testing.T, token string, expectedStatusCode int) []string {
	status, body := readWithRetry(t, func() *http.Response {
		return clientfavorites.HttpGetFavorites(t, token)
	})

	require.Equalf(t, expectedStatusCode, status,
		"HTTP status code must be %d, response: %s", expectedStatusCode, body)

	var ids []string
	for _, item := range gjson.Get(body, "items").Array() {
		id := gjson.Get(item.String(), "advertisement.id").String()
		if id == "" {
			id = gjson.Get(item.String(), "advertisement_id").String()
		}
		if id == "" {
			id = gjson.Get(item.String(), "id").String()
		}
		ids = append(ids, id)
	}
	return ids
}

func GetFavoritesRaw(t *testing.T, token string) (int, string) {
	return readWithRetry(t, func() *http.Response {
		return clientfavorites.HttpGetFavorites(t, token)
	})
}
