package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"net/http"
	"time"
)

func getCallBackJob(request *http.Request, retryPeriod time.Duration, retryCount int) func() {

	return func() {
		doRequest(request, retryPeriod, retryCount)
	}

}

func doRequest(request *http.Request, retryPeriod time.Duration, retryLeft int) {

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
			doRequest(request, retryPeriod, retryLeft-1)
		}()
	} else if fluResp.HttpStatusCode == http.StatusOK {
		go putDbLog(completedFLUs, SUCCESS, *fluResp)

	} else {
		go putDbLog(completedFLUs, "ERROR", *fluResp)
	}
}
