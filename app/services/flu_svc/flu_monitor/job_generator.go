package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/models/status_codes"
	"net/http"
)

func callBackJob(request http.Request){
	client := &http.Client{}
	resp, err := client.Do(&request)
	if err != nil {
		plog.Error("HTTP Error:", err)
		return
	}

	fluResp, status := validationErrorCallback(resp)
	fluResp.FluStatusCode = status

	if status == status_codes.Success {

		go putDbLog(completedFLUs, SUCCESS, *fluResp)

	} else if status == status_codes.CallBackFailure && shouldRetryHttp(projectId) {
		//not successful scenarios
		retryIdsList = append(retryIdsList, projectId)

	} else {
		completedFLUs := deleteFromFeedLinePipe(projectId, fluOutObj)
		go putDbLog(completedFLUs, "ERROR", *fluResp)
	}
	return
}