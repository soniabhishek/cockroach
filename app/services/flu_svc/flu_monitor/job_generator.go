package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/models/status_codes"
	"net/http"
)
func getCallBackJob(request *http.Request, retryCount int, retryPeriod int) func(){

	return func (){

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

		} else if status == status_codes.CallBackFailure && retryCount>0 {
			//not successful scenarios
			//callBackJob(request, retryCount--, retryPeriod)

		} else {
			go putDbLog(completedFLUs, "ERROR", *fluResp)
		}
		return
	}

}
