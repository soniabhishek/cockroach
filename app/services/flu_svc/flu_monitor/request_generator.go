package flu_monitor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
	"net/http"
)

var requestGenPoolCount = make(map[uuid.UUID]int) // Hash map to store queues

func checkRequestGenPool(projectConfig projectLookup) {
	limit := projectConfig.maxFluCount
	queue := queues[projectConfig.projectId]

	//TODO make the number of pools configurable.
	if requestGenPoolCount[projectConfig.projectId] < 1 {
		requestGenPoolCount[projectConfig.projectId]++
		for {
			var fluOutputObj []fluOutputStruct
			//TODO add wait time restriction may be. in case inbound flu rate is very less.
			for i := limit - 1; i >= 0; i-- {
				receiver := queue.Receiver()

				flu := <-receiver

				defer flu.ConfirmReceive()

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
			plog.Info("SENDING FLUs COUNT: ", limit)
			req, err := createRequest(projectConfig.config, fluOutputObj)
			if err != nil {
				plog.Error("Error while creating request", err, " fluOutputObj : ", fluOutputObj)
			}

			job := bulk_processor.NewJob(getCallBackJob(&req, defaultRetryTimePeriod, defaultRetryCount))
			projectConfig.jobManager.PushJob(job)
		}
	}
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
