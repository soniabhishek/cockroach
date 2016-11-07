package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	//	"io/ioutil"
	"net/http"

	"fmt"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"io/ioutil"
)

var luigiBaseApiUrl = config.LUIGI_API.Get()
var encryptionUrl = luigiBaseApiUrl + "/api/scale"

func GetLuigiClient() *luigiClient {
	return &luigiClient{}
}

type luigiClient struct {
}

type encryptionResponse struct {
	Images []models.JsonF `json:"images"`
}

type encryptionRequest struct {
	Image_urls interface{} `json:"image_urls"`
}

func (*luigiClient) GetEncryptedUrls(images interface{}) ([]models.JsonF, error) {

	plog.Debug("GetEnc", images.(string))
	bty, _ := json.Marshal(encryptionRequest{images.(string)})

	req, _ := http.NewRequest("POST", encryptionUrl, bytes.NewBuffer(bty))
	req.Header.Add("content-type", "application/json")
	plog.Debug("2", fmt.Sprintf("%s", bty))

	res, err := http.DefaultClient.Do(req)
	plog.Debug("3")

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {

		return nil, errors.New("Error occured in luigi")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var encResponse encryptionResponse
	//err = json.NewDecoder(encResponse).Decode(res.Body)
	plog.Debug("3")
	err = json.Unmarshal(body, &encResponse)
	if err != nil {
		return []models.JsonF{}, err
	}
	return encResponse.Images, nil
}
