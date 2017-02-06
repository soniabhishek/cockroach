package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/call_back_unit_pipe"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
	"net/http"
)

func getCallBackJob(pHandler *ProjectHandler, cbu *call_back_unit_pipe.CBU) func() {
	plog.Trace("FluMonitor", "Getting Callback job", cbu.FlusSent)
	return func() {
		req, err := createRequest(cbu.ProjectConfig, cbu.FluOutputObj)
		if err != nil {
			plog.Error("FluMonitor", err, plog.M("Error while creating request"), plog.MP(log_tags.POSTBACK_REQUEST, cbu.FluOutputObj))
		}
		client := http.DefaultClient
		resp, err := client.Do(&req)
		if err != nil {
			plog.Error("HTTP Error:", err)
			putDbLogCustom(cbu.FlusSent, "Error", models.JsonF{"HTTP Error": err.Error()})
			cbu.ConfirmReceive()
			return
		}

		fluResp := ParseFluResponse(resp)

		plog.Info("Flu monitor", "Response received", *fluResp)

		//validFlus, invalidFLus := getFlusStatus(cbu.FlusSent, fluResp.Invalid_Flus)

		putDbLog(cbu.FlusSent, *fluResp)

		if shouldRetry(fluResp, cbu.RetryLeft) {
			cbu.RetryLeft--
			pHandler.retryQueue.Push(*cbu)
		}
		cbu.ConfirmReceive()

		plog.Info("FluMonitor", "Job Executed", "ProjectId: "+pHandler.projectId.String(), "FluIDs: ", getFluIds(cbu.FlusSent))
	}
}
