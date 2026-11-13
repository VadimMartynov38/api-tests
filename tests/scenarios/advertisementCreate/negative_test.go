package advertisementCreate

import (
	"net/http"
	"strings"

	"github.com/stretchr/testify/require"

	clientadv "api-tests-template/internal/client/http/advertisements"
	adv "api-tests-template/internal/managers/advertisements"
)

func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	cases := []struct {
		name        string
		title       string
		description string
		price       int
		hasImage    bool
	}{
		{"empty title", "", "desc", 100, true},
		{"negative price", "Test ad", "desc", -1, true},
		{"too long title", strings.Repeat("a", 1000), "desc", 100, true},
		{"missing image", "Test ad", "desc", 100, false},
	}
	for _, tc := range cases {
		tc := tc
		s.Run(tc.name, func() {
			status, body := adv.GetRawResponse(s.T(), s.loginData.Token, clientadv.AdvertisementRequest{
				Title:       tc.title,
				Description: tc.description,
				Price:       tc.price,
				HasImage:    tc.hasImage,
			})
			require.Equal(s.T(), http.StatusBadRequest, status,
				"expected 400, got: %d, body: %s", status, body)
			require.NotEmpty(s.T(), body, "error message is empty")
		})
	}
}

func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name  string
		token string
	}{
		{"no token", ""},
		{"invalid token", "invalid_token_123"},
	}
	for _, tc := range cases {
		tc := tc
		s.Run(tc.name, func() {
			status, body := adv.GetRawResponse(s.T(), tc.token, clientadv.AdvertisementRequest{
				Title:       "Test",
				Description: "Desc",
				Price:       100,
				HasImage:    true,
			})
			require.Equal(s.T(), http.StatusUnauthorized, status,
				"expected 401, got: %d, body: %s", status, body)
			require.NotEmpty(s.T(), body, "error message is empty")
		})
	}
}
