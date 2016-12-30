package flu_monitor

import (
	"bytes"
	"encoding/json"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"github.com/pkg/errors"
	"net/http"
)

type invalidFlu struct {
	Flu_Id  string `json:"flu_id"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func createRequest(config models.ProjectConfiguration, fluProjectResp []models.FluOutputStruct) (request http.Request, err error) {

	//TODO change someshit
	if len(fluProjectResp) < 1 {
		return request, errors.New("someshit")
	}

	plog.Info("Flu output", "sendBackToClient", config.ProjectId)

	url := config.PostBackUrl
	//url := "http://localhost:8080/JServer/HelloServlet"
	plog.Trace("URL:>", url, "|ID: ", config.ProjectId, "|Body:", fluProjectResp)

	sendResp := make(map[string][]models.FluOutputStruct)
	sendResp["feed_line_units"] = fluProjectResp
	jsonBytes, err := json.Marshal(sendResp)
	if err != nil {
		plog.Error("JSON Marshalling Error:", err)
		return request, err
	}
	jsonBytes = utilities.ReplaceEscapeCharacters(jsonBytes)
	plog.Trace("Sending JSON:", string(jsonBytes))

	//fmt.Println(hex.EncodeToString(sig.Sum(nil)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	for headerKey, headerVal := range config.Headers {
		req.Header.Set(headerKey, headerVal.(string))

	}
	addSendBackAuth(req, config, jsonBytes)

	return *req, nil
}
