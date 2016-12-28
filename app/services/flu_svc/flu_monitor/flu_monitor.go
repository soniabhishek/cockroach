package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
	"github.com/crowdflux/angel/utilities"
	"sync"
)

type FluMonitor struct {
}

type projectLookup struct {
	projectId      uuid.UUID
	config         models.ProjectConfiguration
	maxFluCount    int
	postBackUrl    string
	queryFrequency int
	jobManager     bulk_processor.JobManager
}

var activeProjectsLookup = make(map[uuid.UUID]projectLookup) // Hash map to store config
var queues = make(map[uuid.UUID]feed_line.Fl)                // Hash map to store queues
var dispatcherStarter sync.Once
var dispatcher = bulk_processor.NewDispatcher(utilities.GetInt(config.MAX_WORKERS.Get()))

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	//TODO rename clientQ to projectQ, present to ok
	projectQ, ok := queues[flu.ProjectId]
	if !ok {
		//TODO toString in config may be
		clientQ := feed_line.New("Gate-Q-" + flu.ProjectId.String())
		queues[flu.ProjectId] = clientQ
	}
	projectQ.Push(feed_line.FLU{FeedLineUnit: flu})

	pConfig := checkProjectConfig(flu)

	checkRequestGenPool(pConfig)

	dispatcherStarter.Do(func() {
		dispatcher.Start()
	})

	return nil
}

func checkProjectConfig(flu models.FeedLineUnit) projectLookup {

	//TODO activeProjects to activeProjectConfigurations
	projectLookup, ok := activeProjectsLookup[flu.ProjectId]
	if !ok {
		fpsRepo := project_configuration_repo.New()
		fpsModel, err := fpsRepo.Get(flu.ProjectId)
		if err != nil {
			plog.Error("Error while getting Project configuratin", err, " ProjectId:", flu.ProjectId)
		}

		// reconsider
		maxFluCount := getMaxFluCount(fpsModel)
		postbackUrl := fpsModel.PostBackUrl
		//TODO Handle invalid url
		queryFrequency := getQueryFrequency(fpsModel)

		jm := bulk_processor.NewJobManager(projectLookup.queryFrequency, flu.ProjectId.String())
		dispatcher.AddJobManager(jm)

		projectLookup = projectLookup{flu.ProjectId, fpsModel, maxFluCount, postbackUrl, queryFrequency, *jm}
		activeProjectsLookup[flu.ProjectId] = projectLookup
	}
	return projectLookup

}
