package advertisements

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

var pngBytes = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
	0x0D, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
}

type AdvertisementRequest struct {
	Title       string
	Description string
	Price       int
	HasImage    bool
}

func HttpPostAdvertisement(t *testing.T, token string, request AdvertisementRequest) *http.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("title", request.Title)
	_ = writer.WriteField("description", request.Description)
	_ = writer.WriteField("price", fmt.Sprint(request.Price))
	if request.HasImage {
		part, _ := writer.CreateFormFile("photos", "img1.png")
		_, _ = part.Write(pngBytes)
	}
	_ = writer.Close()

	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().Response
}

func HttpGetAdvertisement(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", adID).
		Expect(t).
		End().Response
}

func HttpGetAdvertisements(t *testing.T, token, text string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		Query("search", text).
		Expect(t).
		End().Response
}
