package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"io/ioutil"
	"net/http"
)

var luigiBaseApiUrl = config.LUIGI_API.Get()
var encryptionUrl = luigiBaseApiUrl + "/api/scale"

func GetLuigiClient() *luigiClient {
	return &luigiClient{}
}

type luigiClient struct {
}

type encryptionResponse struct {
	Images map[string]LuigiResponse `json:"images"`
}

type encryptionRequest struct {
	Image_urls []string `json:"image_urls"`
}

type LuigiResponse struct {
	Value        bool           `json:"valid"`
	SourceUrl    string         `json:"source_url"`
	PlaymentUrl  string         `json:"playment_url"`
	ScaledImages []models.JsonF `json:"scaled_images"`
}

func (*luigiClient) GetEncryptedUrls(images []string) (map[string]LuigiResponse, error) {
	bty, _ := json.Marshal(encryptionRequest{images})

	req, _ := http.NewRequest("POST", encryptionUrl, bytes.NewBuffer(bty))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {

		return nil, errors.New("Error occured in luigi. Check image urls")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var encResponse encryptionResponse
	err = json.Unmarshal(body, &encResponse)

	if err != nil {
		return map[string]LuigiResponse{}, err
	}
	return encResponse.Images, nil
}
