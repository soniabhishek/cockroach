package flu_monitor

import (
	"bytes"
	"encoding/json"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/status_codes"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"net/http"
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

func sendBackToClient(config models.ProjectConfiguration, fluProjectResp []fluOutputStruct) (*FluResponse, status_codes.StatusCode) {

	if len(fluProjectResp) < 1 {
		return &FluResponse{}, status_codes.NoFluToSend
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
		return &FluResponse{}, status_codes.UnknownFailure
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

	job:=Job{Request:*req}
	JobQueue<-job

	 //separate
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		plog.Error("HTTP Error:", err)
		return &FluResponse{}, status_codes.UnknownFailure
	}

	fluResp, status := validationErrorCallback(resp)
	fluResp.FluStatusCode = status
	return fluResp, status
}
