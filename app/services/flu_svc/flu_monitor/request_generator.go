package flu_monitor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/http_request_unit_pipe"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

var requestGenPoolCount = make(map[uuid.UUID]int) // Hash map to store queues

func checkRequestGenPool(projectConfig projectHandler) {
	limit := projectConfig.maxFluCount
	queue := projectConfig.queue

	//TODO make the number of pools configurable.
	if requestGenPoolCount[projectConfig.projectId] < 1 {
		requestGenPoolCount[projectConfig.projectId]++
		for {
			var fluOutputObj []models.FluOutputStruct
			var flusSent = make(map[uuid.UUID]feed_line.FLU)
			//TODO add wait time restriction may be. in case inbound flu rate is very less.
			for i := limit - 1; i >= 0; i-- {
				receiver := queue.Receiver()

				select {
				case flu := <-receiver:
					defer flu.ConfirmReceive()

				}
				flu := <-receiver

				defer flu.ConfirmReceive()

				result, ok := flu.Build[RESULT]
				if !ok {
					result = models.JsonF{}

				}

				fluOutputObj = append(fluOutputObj, models.FluOutputStruct{
					ID:          flu.ID,
					ReferenceId: flu.ReferenceId,
					Tag:         flu.Tag,
					Status:      STATUS_COMPLETED,
					Result:      result,
				})
				flusSent[flu.ID] = flu
			}
			plog.Info("SENDING FLUs COUNT: ", limit)
			requestQueues.Push(http_request_pipe.FMCR{FluOutputObj: fluOutputObj})
		}
	}
}

func addSendBackAuth(req *http.Request, fpsModel models.ProjectConfiguration, bodyJsonBytes []byte) {
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
