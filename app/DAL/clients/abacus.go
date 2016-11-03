package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/config"
	"io/ioutil"
	"net/http"
)

var abacusApiUrl = config.ABACUS_API.Get()
var algorithmUrl = abacusApiUrl + "/api/review"

func GetAbacusClient() *abacusClient {
	return &abacusClient{}
}

type abacusClient struct {
}

type algorithmResponse struct {
	prediction string      `json:"prediction"`
	Success    bool        `json:"success"`
	Error      interface{} `json:"error"`
}

type algorithmRequest struct {
	Input string `json:"review"`
}

func (*abacusClient) Predict(text string) (string, error) {

	bty, _ := json.Marshal(algorithmRequest{text})

	req, _ := http.NewRequest("POST", algorithmUrl, bytes.NewBuffer(bty))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {

		return "", errors.New("Error occured in abacus")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var algorithmResponse algorithmResponse

	err = json.Unmarshal(body, &algorithmResponse)
	if err != nil {
		return "", err
	}
	return algorithmResponse.prediction, nil
}
