package flu_monitor

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"sync"
	"time"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
)

type FluMonitor struct {
}

type projectConfig struct {
	projectId      uuid.UUID
	config         models.ProjectConfiguration
	maxFluCount    int
	postBackUrl    string
	queryFrequency int
}

var retryCount = make(map[uuid.UUID]int)
var mutex = &sync.RWMutex{}
var dbLogger = feed_line_repo.NewLogger()

var retryTimePeriod = time.Duration(utilities.GetInt(config.RETRY_TIME_PERIOD.Get())) * time.Millisecond
var fluThresholdDuration = int64(utilities.GetInt(config.FLU_THRESHOLD_DURATION.Get()))
var monitorTimePeriod = time.Duration(utilities.GetInt(config.MONITOR_TIME_PERIOD.Get())) * time.Millisecond
var retryThreshold = utilities.GetInt(config.FLU_RETRY_THRESHOLD.Get())

var activeProjects = make(map[uuid.UUID]projectConfig) // Hash map to store config
var queues = make(map[uuid.UUID]feed_line.Fl)          // Hash map to store queues

var defaultFluThresholdCount = utilities.GetInt(config.DEFAULT_FLU_THRESHOLD_COUNT.Get())

var dispatcherStater sync.Once
var dispatcher = bulk_processor.NewDispatcher(0)

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	//TODO rename clientQ to projectQ, present to ok
	clientQ, present := queues[flu.ProjectId]
	if !present {
		//TODO add prefix (output queue)
		clientQ := feed_line.New(flu.ProjectId.String())
		queues[flu.ProjectId] = clientQ
	}
	clientQ.Push(feed_line.FLU{FeedLineUnit: flu})

	checkProjectConfig(flu)

	dispatcherStater.Do(func() {
		dispatcher.Start()
	})

	//SendData(c)

	req,err:=makeRequest(activeProjects[flu.ProjectId])
	//generateRequest()

	//generateJob()

	//push to bulk_processor

	return nil
}

func checkProjectConfig(flu models.FeedLineUnit) {

	//TODO activeProjects to activeProjectConfigurations
	value, valuePresent := activeProjects[flu.ProjectId]
	if !valuePresent {
		fpsRepo := project_configuration_repo.New()
		fpsModel, err := fpsRepo.Get(flu.ProjectId)
		if err!=nil {
			plog.Error("Error while getting Project configuratin", err, " ProjectId:",flu.ProjectId)
		}

		// reconsider
		maxFluCount := getMaxFluCount(fpsModel)
		postbackUrl := fpsModel.PostBackUrl
		//TODO Handle invalid url
		queryFrequency := getQueryFrequency(fpsModel)
		value = projectConfig{flu.ProjectId, fpsModel, maxFluCount, postbackUrl, queryFrequency}
		activeProjects[flu.ProjectId] = value
	}

	jm:=bulk_processor.NewJobManager(1, flu.ProjectId.String())
	dispatcher.AddJobManager(jm)
}
