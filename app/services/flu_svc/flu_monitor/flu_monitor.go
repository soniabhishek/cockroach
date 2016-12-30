package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
	"github.com/crowdflux/angel/utilities"
	"sync"
)

type FluMonitor struct {
	projectHandlers map[uuid.UUID]projectHandler
	bulkProcessor   *bulk_processor.Dispatcher
}

type projectHandler struct {
	projectId      uuid.UUID
	config         models.ProjectConfiguration
	maxFluCount    int
	postBackUrl    string
	queryFrequency int
	jobManager     bulk_processor.JobManager
	queue          feed_line.Fl
}

func New() *FluMonitor {
	return &FluMonitor{
		projectHandlers: make(map[uuid.UUID]projectHandler),
		bulkProcessor:   bulk_processor.NewDispatcher(services.AtoiOrPanic(config.MAX_WORKERS.Get())),
	}
}

var dispatcherStarter sync.Once

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	projectQ := fm.projectHandlers[flu.ProjectId].queue

	projectQ.Push(feed_line.FLU{FeedLineUnit: flu})

	pHandler := fm.getProjectHandler(flu)

	checkRequestGenPool(pHandler)

	defer dispatcherStarter.Do(func() {
		fm.bulkProcessor.Start()
	})

	job := bulk_processor.NewJob(getCallBackJob(retryTimePeriod, retryCount))
	pHandler.jobManager.PushJob(job)

	return nil
}

func (fm *FluMonitor) getProjectHandler(flu models.FeedLineUnit) projectHandler {

	//TODO activeProjects to activeProjectConfigurations
	projectLookup, ok := fm.projectHandlers[flu.ProjectId]
	if !ok {
		pcRepo := project_configuration_repo.New()
		pc, err := pcRepo.Get(flu.ProjectId)
		if err != nil {
			plog.Error("Error while getting Project configuratin", err, " ProjectId:", flu.ProjectId)
		}

		// reconsider
		maxFluCount := getMaxFluCount(pc)
		postbackUrl := pc.PostBackUrl
		//TODO Handle invalid url
		queryFrequency := getQueryFrequency(pc)

		queue := feed_line.New("Gate-Q-" + flu.ProjectId.String())
		jm := bulk_processor.NewJobManager(projectLookup.queryFrequency, flu.ProjectId.String())
		fm.bulkProcessor.AddJobManager(jm)

		projectLookup = projectHandler{flu.ProjectId, pc, maxFluCount, postbackUrl, queryFrequency, *jm, queue}
		fm.bulkProcessor[flu.ProjectId] = projectLookup
	}
	return projectLookup

}
