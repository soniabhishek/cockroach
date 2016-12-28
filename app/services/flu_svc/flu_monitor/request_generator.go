package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"fmt"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
)

func makeRequest(projectConfig projectConfig) (error) {
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

		select {
		case flu := <-receiver:
			defer flu.ConfirmReceive()
		default:
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

	req, err:=createRequest(projectConfig.config, fluOutputObj)

	job:= bulk_processor.NewJob(getCallBackJob(req))
	job:=Job{Request:*req}
	JobQueue<-job

	return err
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

