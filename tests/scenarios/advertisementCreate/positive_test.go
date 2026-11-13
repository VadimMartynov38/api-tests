package advertisementCreate

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	clientadv "api-tests-template/internal/client/http/advertisements"
	adv "api-tests-template/internal/managers/advertisements"
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

func (s *TestSuite) TestCreateAdvertisementFull() {
	title := "Test ad " + utils.RandomString(10)
	description := "Full description of test ad " + utils.RandomString(10)
	price := 12345

	var adID string
	s.Run("Create advertisement with all parameters and a photo", func() {
		resp := adv.CreateAdvertisement(s.T(), s.loginData.Token, clientadv.AdvertisementRequest{
			Title:       title,
			Description: description,
			Price:       price,
			HasImage:    true,
		}, http.StatusCreated)
		require.NotEmpty(s.T(), resp.ID, "ad id is empty in response")
		adID = resp.ID
	})

	s.Run("Ad is available via GET /advertisement?id= with the same parameters", func() {
		got := adv.GetAdvertisementByID(s.T(), s.loginData.Token, adID, http.StatusOK)
		require.Equal(s.T(), title, got.Title, "title mismatch")
		require.Equal(s.T(), description, got.Description, "description mismatch")
		require.Equal(s.T(), float64(price), got.Price, "price mismatch")
		require.NotEmpty(s.T(), got.Photos, "photo list is empty in the ad object")
	})

	s.Run("Ad is found in search by full title match", func() {
		list := adv.GetAdvertisementsList(s.T(), s.loginData.Token, title)
		require.NotEmpty(s.T(), list.Items, "search returned no results")

		found := false
		for _, item := range list.Items {
			if item.ID == adID {
				found = true
			}
		}
		require.True(s.T(), found, "created ad was not found in search results")
	})

	s.Run("Cleanup: delete the created advertisement", func() {
		adv.DeleteAdvertisement(s.T(), s.loginData.Token, adID, http.StatusNoContent)
	})
}
