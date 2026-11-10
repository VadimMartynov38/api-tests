package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	existingAdTitle = "Selling Fender guitar"
	searchBeginning = "Selling"
	searchMiddle    = "Fender"
	searchEnd       = "guitar"
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

	base.Precondition("Check that search constants are set and authorize")
	require.NotEmpty(s.T(), existingAdTitle, "set existingAdTitle in search_test.go")
	require.NotEmpty(s.T(), searchBeginning, "set searchBeginning in search_test.go")
	require.NotEmpty(s.T(), searchMiddle, "set searchMiddle in search_test.go")
	require.NotEmpty(s.T(), searchEnd, "set searchEnd in search_test.go")

	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	var list advModels.AdvertisementsList
	s.Run("Get search results by middle substring", func() {
		list = advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)
	})

	s.Run("Response has at least one advertisement", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Our advertisement is among results", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"beginning", searchBeginning},
		{"middle", searchMiddle},
		{"end", searchEnd},
		{"lower case", strings.ToLower(searchMiddle)},
		{"upper case", strings.ToUpper(searchMiddle)},
		{"mixed case", utils.MixCase(searchMiddle)},
	}
	for _, tc := range cases {
		tc := tc
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
