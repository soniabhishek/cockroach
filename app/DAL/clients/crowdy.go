package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
)

var crowdyBaseApiUrl = config.Get(config.CROWDY_BASE_API)
var pushFluUrl = crowdyBaseApiUrl + "/crowdsourcing_gateway?action=add_flu"
var authkey = config.Get(config.CROWDY_AUTH_KEY)

func GetCrowdyClient() *crowdyClient {
	return &crowdyClient{}
}

type crowdyClient struct {
}
type pushFluResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}

type pushFluReq struct {
	Flu models.FeedLineUnit `json:"flu"`
}

func (*crowdyClient) PushFLU(flu models.FeedLineUnit) (bool, error) {

	bty, _ := json.Marshal(pushFluReq{flu})

	req, _ := http.NewRequest("POST", pushFluUrl, bytes.NewBuffer(bty))

	req.Header.Add("authorization", authkey)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var pushFluResp pushFluResponse

	err = json.Unmarshal(body, &pushFluResp)
	if err != nil {
		return false, err
	}

	if pushFluResp.Success {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(pushFluResp.Error))
	}
}
