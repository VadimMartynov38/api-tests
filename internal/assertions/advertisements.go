package assertions

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	for _, ad := range items {
		if strings.EqualFold(ad.Title, title) {
			return
		}
	}
	require.Failf(t, "advertisement not found",
		"no advertisement with title %q in result of %d items", title, len(items))
}
