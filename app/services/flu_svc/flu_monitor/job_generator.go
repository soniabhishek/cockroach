package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
	"time"
)

func getCallBackJob(request *http.Request, retryPeriod time.Duration, retryCount int, flusSent map[uuid.UUID]feed_line.FLU) func() {
	return func() {
		doRequest(request, retryPeriod, retryCount, flusSent)
	}

}

func doRequest(request *http.Request, retryPeriod time.Duration, retryLeft int, flusSent map[uuid.UUID]feed_line.FLU) {

	if retryLeft == 0 {
		return
	}
	client := http.DefaultClient
	resp, err := client.Do(request)
	if err != nil {
		plog.Error("HTTP Error:", err)
		return
	}

	fluResp, shouldRetry := shouldRetry(resp)

	if shouldRetry {
		go func() {
			time.Sleep(retryPeriod)
			doRequest(request, retryPeriod, retryLeft-1, flusSent)
		}()
	} else if fluResp.HttpStatusCode == http.StatusOK {
		go putDbLog(getAllFlus(flusSent), SUCCESS, *fluResp)

	} else {
		validFlus, invalidFLus := getFlusStatus(flusSent, fluResp)
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
