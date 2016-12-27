package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"fmt"
	"github.com/crowdflux/angel/app/models/status_codes"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func makeRequest(projectConfig projectConfig) (http.Request, error) {
	// getFluOutputObj(projectConfig)
	limit := projectConfig.maxFluCount
	plog.Info("SENDING FLUs COUNT: ", limit)
	queue := queues[projectConfig.projectId]
	var fluOutputObj []fluOutputStruct
	for i := limit - 1; i >= 0; i-- {
		receiver := queue.Receiver()

		flu := <-receiver

		defer flu.ConfirmReceive() // defer what happens

		// if queue empty, break
		// adjust with the wait time in switch case channel
		// some bufferred channel logic instead of counting in for loop?

		select {
		case flu := <-receiver:
			defer flu.ConfirmReceive()
		default:
			delete(activeProjects, projectConfig.projectId)
			fmt.Println("No value ready, moving on.")
		}
		result, ok := flu.Build[RESULT]
		if !ok {
			result = models.JsonF{}
		}

		fluOutputObj = append(fluOutputObj, fluOutputStruct{
			ID:          flu.ID,
			ReferenceId: flu.ReferenceId,
			Tag:         flu.Tag,
			Status:      STATUS_OK,
			Result:      result,
		})
	}

	// http call and retry logic
	// make request
	// keep retrying in case of failure
	// if success availableQps --
	// defer flu.ConfirmReceive, if the server crashes before the httpcall it stays in queue??

	createRequest(projectConfig.config, fluOutputObj)

	return req,nil
}

func addSendBackAuth(req *http.Request, fpsModel models.ProjectConfiguration, bodyJsonBytes []byte) {
	hmacKey := fpsModel.Options[HMAC_KEY]
	if hmacKey != nil {
		// ToDo add this when encrypted will be in DB
		//hmacKey, _ := utilities.Decrypt(hmacKey.(string))
		sig := hmac.New(sha256.New, []byte(hmacKey.(string)))
		sig.Write([]byte(string(bodyJsonBytes)))
		hmac := hex.EncodeToString(sig.Sum(nil))
		req.Header.Set(HMAC_HEADER_KEY, hmac)
		plog.Trace("HMAC", hmac)
	}
}

