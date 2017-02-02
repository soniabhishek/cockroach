package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/crowdflux/angel/bulk_processor"
	"sync"
)

type FluMonitor struct {
	projectHandlers   map[uuid.UUID]ProjectHandler
	bulkProcessor     *bulk_processor.Dispatcher
	dispatcherStarter sync.Once
	mutex             sync.Mutex
}

func New() *FluMonitor {
	return &FluMonitor{
		projectHandlers: make(map[uuid.UUID]ProjectHandler),
		bulkProcessor:   bulk_processor.NewDispatcher(services.AtoiOrPanic(config.MAX_WORKERS.Get())),
	}
}

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	pHandler := fm.getOrCreateProjectHandler(flu)

	pHandler.queue.Push(feed_line.FLU{FeedLineUnit: flu})

	fm.dispatcherStarter.Do(func() {
		fm.bulkProcessor.Start()
		plog.Info("Flu Monitor", "Bulk Processor started")
	})

	return nil
}

func (fm *FluMonitor) getOrCreateProjectHandler(flu models.FeedLineUnit) ProjectHandler {

	projectHandler, ok := fm.projectHandlers[flu.ProjectId]
	if !ok {
		pcRepo := project_configuration_repo.New()
		pc, err := pcRepo.Get(flu.ProjectId)
		if err != nil {
			plog.Error("Error while getting Project configuration", err, " ProjectId:", flu.ProjectId)
		}

		pHandler := NewProjectHandler(pc)

		fm.bulkProcessor.AddJobManager(pHandler.jobManager)

		fm.mutex.Lock()
		fm.projectHandlers[flu.ProjectId] = pHandler
		fm.mutex.Unlock()

		go pHandler.startFeedLineProcessor()
		go pHandler.startCBUProcessor()

		projectHandler = pHandler
		plog.Info("Flu Monitor", "Project Handler Set", projectHandler.config.ProjectId)
	}
	return projectHandler
}
