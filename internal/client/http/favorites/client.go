package favorites

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpPostFavorite(t *testing.T, token, adID string) *http.Response {
	payload, err := json.Marshal(map[string]string{"advertisement_id": adID})
	require.NoError(t, err)

	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		Body(string(payload)).
		Expect(t).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.FavoritesPath+"/"+adID).
		Expect(t).
		End().Response
}

func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Expect(t).
		End().Response
}
