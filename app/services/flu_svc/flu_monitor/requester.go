package flu_monitor

import (
	"bytes"
	"encoding/json"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"net/http"
	"github.com/pkg/errors"
)

type fluOutputStruct struct {
	ID          uuid.UUID   `json:"flu_id"`
	ReferenceId string      `json:"reference_id"`
	Tag         string      `json:"tag"`
	Status      string      `json:"status"`
	Result      interface{} `json:"result"`
}

type invalidFlu struct {
	Flu_Id  string `json:"flu_id"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func createRequest(config models.ProjectConfiguration, fluProjectResp []fluOutputStruct)  (http.Request, error){

	//TODO change someshit
	if len(fluProjectResp) < 1 {
		return nil, errors.New("someshit")
	}

	plog.Info("Flu output", "sendBackToClient", config.ProjectId)

	url := config.PostBackUrl
	//url := "http://localhost:8080/JServer/HelloServlet"
	plog.Trace("URL:>", url, "|ID: ", config.ProjectId, "|Body:", fluProjectResp)

	sendResp := make(map[string][]fluOutputStruct)
	sendResp["feed_line_units"] = fluProjectResp
	jsonBytes, err := json.Marshal(sendResp)
	if err != nil {
		plog.Error("JSON Marshalling Error:", err)
		return nil, err
	}
	jsonBytes = utilities.ReplaceEscapeCharacters(jsonBytes)
	plog.Trace("Sending JSON:", string(jsonBytes))

	//fmt.Println(hex.EncodeToString(sig.Sum(nil)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set(CONTENT_TYPE, TYPE_JSON)

	for headerKey, headerVal := range config.Headers {
		req.Header.Set(headerKey, headerVal.(string))

	}
	addSendBackAuth(req, config, jsonBytes)

	return *req, nil

	return fluResp, status
}
