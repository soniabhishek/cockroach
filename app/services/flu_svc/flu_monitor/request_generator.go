package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
)

func makeRequest(projectConfig config) error {
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
		if flu == nil {
			delete(projectConfig, projectConfig.projectId)
			break
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

	go sendBackToClient(projectConfig, fluOutputObj)

	return
}
