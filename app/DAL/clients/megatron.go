package clients

import (
	"bytes"
	"encoding/json"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"io/ioutil"
	"net/http"
)

var megatronApiUrl = config.MEGATRON_API.Get()
var transformationUrl = megatronApiUrl + "/transform"

func GetMegatronClient() *megatronClient {
	return &megatronClient{}
}

type megatronClient struct {
}

type transformationResponse struct {
	Output   models.JsonF `json:"output"`
	Warnings []string     `json:"warnings"`
	Success  bool         `json:"success"`
	Error    interface{}  `json:"error"`
}

type transformationRequest struct {
	TemplateId string       `json:"template_id"`
	Input      models.JsonF `json:"input"`
}

func (*megatronClient) Transform(input models.JsonF, templateId string) (models.JsonF, error) {

	bty, _ := json.Marshal(transformationRequest{templateId, input})

	req, _ := http.NewRequest("POST", transformationUrl, bytes.NewBuffer(bty))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK || err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var transformResp transformationResponse

	err = json.Unmarshal(body, &transformResp)
	if err != nil {
		return models.JsonF{}, err
	}
	return transformResp.Output, nil
}
