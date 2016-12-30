package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/http_request_unit_pipe"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
	"time"
)

var retryQueues = http_request_pipe.New("Retry-Q")     // Hash map to store queues
var requestQueues = http_request_pipe.New("Request-Q") // Hash map to store queues

func getCallBackJob(retryPeriod time.Duration, retryCount int) func() {
	return func() { doRequest(retryPeriod, retryCount) }

}

func doRequest(retryPeriod time.Duration, retryLeft int) {

	requestReceiver := requestQueues.Receiver()
	retryReceiver := retryQueues.Receiver()

	var temp_req http_request_pipe.FMCR
	select {
	case req := <-requestReceiver:
		temp_req = req
		defer req.ConfirmReceive()
	case req := <-retryReceiver:
		temp_req = req
		defer req.ConfirmReceive()

	}

	req, err := createRequest(temp_req.ProjectConfig, temp_req.FluOutputObj)
	if err != nil {
		plog.Error("Error while creating request", err, " fluOutputObj : ", temp_req.FluOutputObj)
	}
	if retryLeft == 0 {
		return
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
			retryQueues.Push(http_request_pipe.FMCR{FluOutputObj: temp_req.FluOutputObj, FlusSent: temp_req.FlusSent, ProjectConfig: temp_req.ProjectConfig, RetryCount: retryLeft - 1})
		}()
	} else if fluResp.HttpStatusCode == http.StatusOK {
		go putDbLog(getAllFlus(temp_req.FlusSent), SUCCESS, *fluResp)

	} else {
		validFlus, invalidFLus := getFlusStatus(temp_req.FlusSent, fluResp)
		go func() {
			putDbLog(invalidFLus, "ERROR", *fluResp)
			putDbLog(validFlus, "SUCCESS", *fluResp)
		}()
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
