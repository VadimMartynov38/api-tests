package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	clientadv "api-tests-template/internal/client/http/advertisements"
	adv "api-tests-template/internal/managers/advertisements"
	fav "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
)

const nonExistentAdID = "11111111-2222-3333-4444-555555555555"

func (s *TestSuite) TestFavoritesNonExistentAd() {
	s.Run("Add a non-existent ad to favorites", func() {
		body := fav.AddFavorite(s.T(), s.loginData.Token, nonExistentAdID, http.StatusNotFound)
		require.NotEmpty(s.T(), body, "error message is empty")
	})

	s.Run("Remove a non-existent ad from favorites", func() {
		fav.DeleteFavorite(s.T(), s.loginData.Token, nonExistentAdID, http.StatusNotFound)
	})
}

func (s *TestSuite) TestFavoritesUnauthorized() {
	invalidToken := "invalid_token_123"

	s.Run("GET /favorites with an invalid token returns 401", func() {
		status, body := fav.GetFavoritesRaw(s.T(), invalidToken)
		require.Equal(s.T(), http.StatusUnauthorized, status, "response: %s", body)
	})

	s.Run("POST /favorites with an invalid token returns 401", func() {
		body := fav.AddFavorite(s.T(), invalidToken, nonExistentAdID, http.StatusUnauthorized)
		require.NotEmpty(s.T(), body, "error message is empty")
	})

	s.Run("DELETE /favorites with an invalid token returns 401", func() {
		body := fav.DeleteFavorite(s.T(), invalidToken, nonExistentAdID, http.StatusUnauthorized)
		require.NotEmpty(s.T(), body, "error message is empty")
	})
}

func (s *TestSuite) TestOwnAdvertisementForbidden() {
	var adID string
	s.Run("Create an own advertisement", func() {
		resp := adv.CreateAdvertisement(s.T(), s.loginData.Token, clientadv.AdvertisementRequest{
			Title:       "Own ad " + utils.RandomString(8),
			Description: "Own ad description",
			Price:       700,
			HasImage:    true,
		}, http.StatusCreated)
		adID = resp.ID
	})

	s.Run("Adding own advertisement to favorites returns 403", func() {
		body := fav.AddFavorite(s.T(), s.loginData.Token, adID, http.StatusForbidden)
		require.Contains(s.T(), body, "Cannot add own advertisement",
			"unexpected error: %s", body)
	})

	s.Run("Cleanup: delete the created advertisement", func() {
		adv.DeleteAdvertisement(s.T(), s.loginData.Token, adID, http.StatusNoContent)
	})
}
