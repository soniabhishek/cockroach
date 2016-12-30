package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/http_request_unit_pipe"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
	"net/http"
	"time"
)

var retryQueues = http_request_pipe.New("Retry-Q")     // Hash map to store queues
var requestQueues = http_request_pipe.New("Request-Q") // Hash map to store queues
var jobGenPoolCount = make(map[uuid.UUID]int)          // Hash map to store queues

func generateJobs(pHandler projectHandler) {
	if jobGenPoolCount[pHandler.projectId] > 0 {
		return
	}

	requestReceiver := requestQueues.Receiver()
	retryReceiver := retryQueues.Receiver()

	for {
		var temp_req http_request_pipe.FMCR
		select {
		case req := <-requestReceiver:
			temp_req = req
			defer req.ConfirmReceive()
		case req := <-retryReceiver:
			temp_req = req
			defer req.ConfirmReceive()

		}
		for {
			job := bulk_processor.NewJob(getCallBackJob(&temp_req, retryTimePeriod, retryCount))
			pHandler.jobManager.PushJob(job)
		}
	}
}

func getCallBackJob(fmcr *http_request_pipe.FMCR, retryPeriod time.Duration, retryLeft int) func() {
	return func() {

		req, err := createRequest(fmcr.ProjectConfig, fmcr.FluOutputObj)
		if err != nil {
			plog.Error("Error while creating request", err, " fluOutputObj : ", fmcr.FluOutputObj)
		}

		client := http.DefaultClient
		resp, err := client.Do(&req)
		if err != nil {
			plog.Error("HTTP Error:", err)
			return
		}

		fluResp, shouldRetry := shouldRetry(resp, retryLeft)

		if shouldRetry {
			go func() {
				time.Sleep(retryPeriod)
				retryQueues.Push(http_request_pipe.FMCR{FluOutputObj: fmcr.FluOutputObj, FlusSent: fmcr.FlusSent, ProjectConfig: temp_req.ProjectConfig, RetryCount: retryLeft - 1})
			}()
		} else if fluResp.HttpStatusCode == http.StatusOK {
			go putDbLog(getAllFlus(fmcr.FlusSent), "SUCCESS", *fluResp)

		} else {
			validFlus, invalidFLus := getFlusStatus(fmcr.FlusSent, fluResp)
			go func() {
				putDbLog(invalidFLus, "ERROR", *fluResp)
				putDbLog(validFlus, "SUCCESS", *fluResp)
			}()
		}
	}
}

func getFlusStatus(flusSent map[uuid.UUID]feed_line.FLU, resp *WebhookResponse) ([]feed_line.FLU, []feed_line.FLU) {
	inFlus := resp.Invalid_Flus
	invalidFlus := make([]feed_line.FLU, len(inFlus))
	validFlus := make([]feed_line.FLU, len(flusSent)-len(inFlus))

	for key, inFlu := range inFlus {
		flu_id := uuid.FromStringOrNil(inFlu.Flu_Id)
		flu, ok := flusSent[flu_id]
		if ok {
			invalidFlus[key] = flu
			delete(flusSent, flu_id)
		}
	}

	i := 0
	for _, flu := range flusSent {
		validFlus[i] = flu
		i++
	}

	return validFlus, invalidFlus
}

func getAllFlus(flusSent map[uuid.UUID]feed_line.FLU) []feed_line.FLU {
	validFlus := make([]feed_line.FLU, len(flusSent))

	i := 0
	for _, flu := range flusSent {
		validFlus[i] = flu
		i++
	}

	return validFlus
}
