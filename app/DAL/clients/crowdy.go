package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"io/ioutil"
	"net/http"
)

var crowdyBaseApiUrl = config.Get(config.CROWDY_BASE_API)
var pushFluUrl = "http://localhost:9000/api/crowdsourcing_gateway?action=add_flu"

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

	req.Header.Add("authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6ImQyM2ZlOWQ5LTJiOTgtNGE5NS1hM2JkLWZmMDZiMWNkYmZlNCIsImlhdCI6MTQ2Mjk3OTM0MX0.1nvmr2O5KU2RSOflMQNJiCqxHdNlKitzMgBH7JM2ktM")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var pushFluResp pushFluResponse

	err := json.Unmarshal(body, &pushFluResp)
	if err != nil {
		return false, err
	}

	if pushFluResp.Success {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprint(pushFluResp.Error))
	}
}
