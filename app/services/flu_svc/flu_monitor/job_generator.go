package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/models/status_codes"
	"net/http"
	"time"
)
func getCallBackJob(request *http.Request, retryPeriod time.Duration, retryCount int) func(){

	return func (){
		doRequest(request, retryPeriod, retryCount)
	}

}

func doRequest(request *http.Request, retryPeriod time.Duration,retryLeft int){

	if retryLeft==0{
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		plog.Error("HTTP Error:", err)
		return
	}

	fluResp, status := validationErrorCallback(resp)
	fluResp.FluStatusCode = status


	if status == status_codes.Success {

		go putDbLog(completedFLUs, SUCCESS, *fluResp)

	} else if status == status_codes.CallBackFailure {
		go func(){
			time.Sleep(retryPeriod)
			doRequest(request, retryPeriod, retryLeft-1)
		}()
	} else {
		go putDbLog(completedFLUs, "ERROR", *fluResp)
	}
}
