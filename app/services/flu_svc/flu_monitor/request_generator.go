package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"fmt"
	"github.com/crowdflux/angel/app/models/status_codes"
	"net/http"
)

func makeRequest(projectConfig projectConfig) error {
	// getFluOutputObj(projectConfig)
	limit := projectConfig.maxFluCount
	plog.Info("SENDING FLUs COUNT: ", limit)
	queue := queues[projectConfig.projectId]
	var fluOutputObj []fluOutputStruct
	for i := limit - 1; i >= 0; i-- {
		receiver := queue.Receiver()

		flu := <-receiver

		defer flu.ConfirmReceive() // defer what happens

		// if queue empty, break
		// adjust with the wait time in switch case channel
		// some bufferred channel logic instead of counting in for loop?

		select {
		case flu := <-receiver:
			defer flu.ConfirmReceive()
		default:
			delete(activeProjects, projectConfig.projectId)
			fmt.Println("No value ready, moving on.")
		}
		result, ok := flu.Build[RESULT]
		if !ok {
			result = models.JsonF{}
		}

		fluOutputObj = append(fluOutputObj, fluOutputStruct{
			ID:          flu.ID,
			ReferenceId: flu.ReferenceId,
			Tag:         flu.Tag,
			Status:      STATUS_OK,
			Result:      result,
		})
	}

	// http call and retry logic
	// make request
	// keep retrying in case of failure
	// if success availableQps --
	// defer flu.ConfirmReceive, if the server crashes before the httpcall it stays in queue??

	sendBackToClient(projectConfig.config, fluOutputObj)

	return nil
}

func (job Job) Do() (*FluResponse, status_codes.StatusCode){
	client := &http.Client{}
	resp, err := client.Do(&job.Request)
	if err != nil {
		plog.Error("HTTP Error:", err)
		return &FluResponse{}, status_codes.UnknownFailure
	}

	fluResp, status := validationErrorCallback(resp)
	fluResp.FluStatusCode = status
	return fluResp, status
}
