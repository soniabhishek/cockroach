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
	"sync"
)

type FluMonitor struct {
	projectHandlers map[uuid.UUID]ProjectHandler
	bulkProcessor   *bulk_processor.Dispatcher
}

type ProjectHandler struct {
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
		projectHandlers: make(map[uuid.UUID]ProjectHandler),
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

	generateJobs(pHandler)

	return nil
}

func (fm *FluMonitor) getProjectHandler(flu models.FeedLineUnit) ProjectHandler {

	projectHandler, ok := fm.projectHandlers[flu.ProjectId]
	if !ok {
		pcRepo := project_configuration_repo.New()
		pc, err := pcRepo.Get(flu.ProjectId)
		if err != nil {
			plog.Error("Error while getting Project configuration", err, " ProjectId:", flu.ProjectId)
		}

		// reconsider
		maxFluCount := getMaxFluCount(pc)
		postbackUrl := pc.PostBackUrl
		queryFrequency := getQueryFrequency(pc)

		queue := feed_line.New("Gate-Q-" + flu.ProjectId.String())
		jm := bulk_processor.NewJobManager(projectHandler.queryFrequency, flu.ProjectId.String())
		fm.bulkProcessor.AddJobManager(jm)

		projectHandler = ProjectHandler{flu.ProjectId, pc, maxFluCount, postbackUrl, queryFrequency, *jm, queue}
		fm.projectHandlers[flu.ProjectId] = projectHandler
	}
	return projectHandler
}
