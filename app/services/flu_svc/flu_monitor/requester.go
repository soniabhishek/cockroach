package flu_monitor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"github.com/pkg/errors"
	"net/http"
)

type invalidFlu struct {
	FluID   string `json:"flu_id"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func createRequest(config models.ProjectConfiguration, fluOutputStructs []models.FluOutputStruct) (request http.Request, err error) {

	//TODO change someshit
	if len(fluOutputStructs) < 1 {
		return request, errors.New("someshit")
	}

	plog.Info("Flu output", "sendBackToClient", config.ProjectId)

	url := config.PostBackUrl
	//url := "http://localhost:8080/JServer/HelloServlet"
	plog.Trace("URL:>", url, "|ID: ", config.ProjectId, "|Body:", fluOutputStructs)

	jsonBytes, err := json.Marshal(struct {
		FeedLineUnits []models.FluOutputStruct `json:"feed_line_units"`
	}{fluOutputStructs})

	if err != nil {
		plog.Error("FluMonitor", err, "JSON Marshalling Error:")
		return request, err
	}
	jsonBytes = utilities.ReplaceEscapeCharacters(jsonBytes)
	plog.Trace("Sending JSON:", string(jsonBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	for headerKey, headerVal := range config.Headers {
		req.Header.Set(headerKey, headerVal.(string))

	}
	if isHmacEnabled(config) {
		addSendBackAuth(req, config, jsonBytes)
	}

	return *req, nil
}

func addSendBackAuth(req *http.Request, fpsModel models.ProjectConfiguration, bodyJsonBytes []byte) {

	//TODO change this to file config instead of db
	hmacKey := getHmacKey(fpsModel)
	hmacHeader := getHmacHeader(fpsModel)

	// ToDo add this when encrypted will be in DB
	//hmacKey, _ := utilities.Decrypt(hmacKey.(string))
	sig := hmac.New(sha256.New, []byte(hmacKey))
	sig.Write([]byte(string(bodyJsonBytes)))
	hmac := hex.EncodeToString(sig.Sum(nil))
	req.Header.Set(hmacHeader, hmac)
	plog.Trace("HMAC", hmac)
}
