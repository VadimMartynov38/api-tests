package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	clientadv "api-tests-template/internal/client/http/advertisements"
	adv "api-tests-template/internal/managers/advertisements"
	fav "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Authorize user with credentials from environment variables")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestFavoritesLifecycle() {
	var foreignAdID string
	var ownUserID string

	s.Run("Create an own advertisement to learn own user_id", func() {
		resp := adv.CreateAdvertisement(s.T(), s.loginData.Token, clientadv.AdvertisementRequest{
			Title:       "Fav probe " + utils.RandomString(8),
			Description: "Fav probe description",
			Price:       500,
			HasImage:    true,
		}, http.StatusCreated)
		require.NotEmpty(s.T(), resp.ID)

		ad := adv.GetAdvertisementByID(s.T(), s.loginData.Token, resp.ID, http.StatusOK)
		ownUserID = ad.UserID
		require.NotEmpty(s.T(), ownUserID, "own user_id is empty")

		adv.DeleteAdvertisement(s.T(), s.loginData.Token, resp.ID, http.StatusNoContent)
	})

	s.Run("Find a foreign advertisement via the public ads list", func() {
		list := adv.GetAdvertisementsList(s.T(), s.loginData.Token, "")
		for _, item := range list.Items {
			if item.UserID != ownUserID && item.ID != "" {
				foreignAdID = item.ID
				break
			}
		}
		require.NotEmpty(s.T(), foreignAdID,
			"no foreign ad found on the board: publish an ad from another account and rerun")
	})

	s.Run("Add the foreign advertisement to favorites", func() {
		status, body := fav.AddFavoriteRaw(s.T(), s.loginData.Token, foreignAdID)
		require.Contains(s.T(), []int{http.StatusOK, http.StatusCreated}, status,
			"expected 200 or 201, response: %s", body)
	})

	s.Run("Advertisement appears in the favorites list", func() {
		favorites := fav.GetFavorites(s.T(), s.loginData.Token, http.StatusOK)
		require.Contains(s.T(), favorites, foreignAdID, "ad is not in favorites")
	})

	s.Run("Remove the advertisement from favorites", func() {
		status, body := fav.DeleteFavoriteRaw(s.T(), s.loginData.Token, foreignAdID)
		require.Contains(s.T(), []int{http.StatusOK, http.StatusNoContent}, status,
			"expected 200 or 204, response: %s", body)
	})

	s.Run("Removed advertisement disappears from the list", func() {
		favorites := fav.GetFavorites(s.T(), s.loginData.Token, http.StatusOK)
		require.NotContains(s.T(), favorites, foreignAdID, "ad is still in favorites")
	})
}
