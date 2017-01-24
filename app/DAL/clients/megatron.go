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

var megatronApiUrl = config.MEGATRON_API.Get()
var transformationUrl = megatronApiUrl + "/transform"
var validationUrl = megatronApiUrl + "/validate"

func GetMegatronClient() *megatronClient {
	return &megatronClient{}
}

type megatronClient struct {
}

type transformationResponse struct {
	Output   models.JsonF  `json:"output"`
	Warnings []interface{} `json:"warnings"`
	Success  bool          `json:"success"`
	Error    interface{}   `json:"error"`
}

type validationResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}

type transformationRequest struct {
	TemplateId string       `json:"template_id"`
	Input      models.JsonF `json:"input"`
}

type validationRequest struct {
	TemplateId string       `json:"template_id"`
	Input      models.JsonF `json:"input"`
}

func (*megatronClient) Transform(input models.JsonF, templateId string) (models.JsonF, error) {

	bty, _ := json.Marshal(transformationRequest{templateId, input})

	req, _ := http.NewRequest("POST", transformationUrl, bytes.NewBuffer(bty))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		return nil, errors.New("Error occured in megatron")
	}

	body, _ := ioutil.ReadAll(res.Body)

	var transformResp transformationResponse

	err = json.Unmarshal(body, &transformResp)
	if err != nil {
		return models.JsonF{}, err
	}
	return transformResp.Output, nil
}

func (*megatronClient) Validate(input models.JsonF, templateId string) (bool, error) {

	bty, _ := json.Marshal(validationRequest{templateId, input})

	req, _ := http.NewRequest("POST", validationUrl, bytes.NewBuffer(bty))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		return false, errors.New("Error occured in megatron")
	}

	body, _ := ioutil.ReadAll(res.Body)

	var validationResp validationResponse

	err = json.Unmarshal(body, &validationResp)
	if err != nil {
		return false, err
	}
	if validationResp.Error != nil {
		return validationResp.Success, errors.New("Error occured in megatron")
	}
	return validationResp.Success, nil
}
