package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/call_back_unit_pipe"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/bulk_processor"
	"github.com/pkg/errors"
	"time"
)

type ProjectHandler struct {
	projectId      uuid.UUID
	config         models.ProjectConfiguration
	maxFluCount    int
	postBackUrl    string
	queryFrequency int
	jobManager     *bulk_processor.JobManager
	queue          feed_line.Fl
	requestQueue   call_back_unit_pipe.CbuQ
	retryQueue     call_back_unit_pipe.CbuQ
}

func NewProjectHandler(pc models.ProjectConfiguration) ProjectHandler {
	maxFluCount := getMaxFluCount(pc)
	postBackUrl := pc.PostBackUrl
	queryFrequency := getQueryFrequency(pc)

	queue := feed_line.New("Gate-Q-" + pc.ProjectId.String())
	jm := bulk_processor.NewJobManager(queryFrequency, pc.ProjectId.String())

	requestQueue := call_back_unit_pipe.New("Request-Q-" + pc.ProjectId.String())
	retryQueue := call_back_unit_pipe.New("Retry-Q-" + pc.ProjectId.String())

	return ProjectHandler{pc.ProjectId, pc, maxFluCount, postBackUrl, queryFrequency, jm, queue, requestQueue, retryQueue}
}

func (pHandler *ProjectHandler) startCBUProcessor() {

	requestQueueReceiver := pHandler.requestQueue.Receiver()
	retryQueueReceiver := pHandler.retryQueue.Receiver()

	for {
		select {
		case cbu := <-requestQueueReceiver:
			job := bulk_processor.NewJob(getCallBackJob(pHandler, &cbu))
			pHandler.jobManager.PushJob(job)
			plog.Info("FluMonitor", "Job Pushed", "RequestQueue", "ProjectId: "+pHandler.projectId.String(), " FluIDs: ", getFluIds(cbu.FlusSent))

		case retryCbu := <-retryQueueReceiver:
			go func(cbu call_back_unit_pipe.CBU) {

				if cbu.RetryLeft >= 0 {
					delayFluRetry(cbu.RetryLeft)
					job := bulk_processor.NewJob(getCallBackJob(pHandler, &cbu))
					pHandler.jobManager.PushJob(job)
					plog.Info("FluMonitor", "Job Pushed", "RetryQueue",
						"ProjectId: "+pHandler.projectId.String(), " FluIDs: ", getFluIds(cbu.FlusSent), "RetryCount: ", cbu.RetryLeft)

				} else {
					plog.Error("FluMonitor", errors.New("0 Retries left but pushed to retry queue"))
				}

			}(retryCbu)
		}
	}
}

func (pHandler *ProjectHandler) startFeedLineProcessor() {
	receiver := pHandler.queue.Receiver()

	for {
		cbu := call_back_unit_pipe.CBU{FlusSent: make(map[uuid.UUID]feed_line.FLU), ProjectConfig: pHandler.config, RetryLeft: MAX_RETRY_COUNT}
		var timer <-chan time.Time

		flu := <-receiver
		addFluToCbu(flu, &cbu)

		timer = time.After(time.Duration(1000/pHandler.queryFrequency) * time.Millisecond)

	OutputObjectGeneratorLoop:
		for i := pHandler.maxFluCount - 2; i >= 0; i-- {

			select {
			case flu := <-receiver:
				addFluToCbu(flu, &cbu)
			case <-timer:
				break OutputObjectGeneratorLoop
			}
		}
		plog.Info("Flu Monitor Project_Handler", "Push to : "+pHandler.projectId.String()+" FLuCount: ", len(cbu.FluOutputObj), " FluIds: ", getFluIds(cbu.FlusSent))
		pHandler.requestQueue.Push(cbu)

		for _, flu := range cbu.FlusSent {
			flu.ConfirmReceive()
		}
	}
}

func addFluToCbu(flu feed_line.FLU, cbu *call_back_unit_pipe.CBU) {
	result, ok := flu.Build[RESULT]
	if !ok {
		result = models.JsonF{}
	}

	cbu.FluOutputObj = append(cbu.FluOutputObj, models.FluOutputStruct{
		ID:          flu.ID,
		ReferenceId: flu.ReferenceId,
		Tag:         flu.Tag,
		Status:      STATUS_COMPLETED,
		Result:      result,
	})
	cbu.FlusSent[flu.ID] = flu
}
func delayFluRetry(retryCount int) {
	time.Sleep(DEFAULT_RETRY_DELAY)
}

func getFluIds(fluMap map[uuid.UUID]feed_line.FLU) (fluIdsString []string) {
	fluIdsString = make([]string, len(fluMap))
	i := 0
	for key, _ := range fluMap {
		fluIdsString[i] = key.String()
		i++
	}
	return
}
