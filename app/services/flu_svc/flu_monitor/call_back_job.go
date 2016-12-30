package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/call_back_unit_pipe"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func getCallBackJob(pHandler *ProjectHandler, cbu *call_back_unit_pipe.CBU, retryCountLeft int) func() {
	return func() {

		req, err := createRequest(cbu.ProjectConfig, cbu.FluOutputObj)
		if err != nil {
			plog.Error("FluMonitor", err, "Error while creating request", " fluOutputObj : ", cbu.FluOutputObj)
		}

		client := http.DefaultClient
		resp, err := client.Do(&req)
		if err != nil {
			plog.Error("HTTP Error:", err)
			return
		}

		fluResp, shouldRetry := shouldRetry(resp, retryCountLeft)

		validFlus, invalidFLus := getFlusStatus(cbu.FlusSent, fluResp.Invalid_Flus)
		putDbLog(invalidFLus, "ERROR", *fluResp)
		putDbLog(validFlus, "SUCCESS", *fluResp)

		if shouldRetry {
			cbu.RetryLeft--
			pHandler.retryQueue.Push(*cbu)
		}
		cbu.ConfirmReceive()
		plog.Info("FluMonitor", "Job Executed", "ProjectId: "+pHandler.projectId.String(), "FluIDs: ", getFluIds(cbu.FlusSent))
	}
}

func getFlusStatus(flusSent map[uuid.UUID]feed_line.FLU, invalidFlus []invalidFlu) ([]feed_line.FLU, []feed_line.FLU) {
	invalidFlusOut := make([]feed_line.FLU, len(invalidFlus))
	validFlus := make([]feed_line.FLU, len(flusSent)-len(invalidFlus))

	for i, invalidFlu := range invalidFlus {
		fluId := uuid.FromStringOrNil(invalidFlu.FluID)
		flu, ok := flusSent[fluId]
		if ok {
			invalidFlusOut[i] = flu
			delete(flusSent, fluId)
		}
	}

	i := 0
	for _, flu := range flusSent {
		validFlus[i] = flu
		i++
	}

	return validFlus, invalidFlusOut
}
