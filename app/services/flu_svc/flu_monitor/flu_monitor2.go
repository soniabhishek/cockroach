package flu_monitor

import (
	"errors"
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
)

type FluMonitor struct {
	PoolIsRunning bool
}

type projectConfig struct {
	projectId      uuid.UUID
	config         models.ProjectConfiguration
	maxFluCount    int
	postBackUrl    string
	queryFrequency int
	queue          feed_line.Fl
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
var totalQps = utilities.GetInt(config.TOTAL_QPS.Get())
var availableQps = totalQps

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	clientQ, present := queues[flu.ProjectId]
	if !present {
		clientQ := feed_line.New(flu.ProjectId.String())
		queues[flu.ProjectId] = clientQ
	}
	clientQ.Push(feed_line.FLU{FeedLineUnit: flu})

	saveProjectConfig(flu)

	fm.servicePoolStart()

	// configure maxworkers
	dispatcher := NewDispatcher(100)
	dispatcher.Run()

	return nil
}

func saveProjectConfig(flu models.FeedLineUnit) error{
	value, valuePresent := activeProjects[flu.ProjectId]
	if !valuePresent {
		fpsRepo := project_configuration_repo.New()
		fpsModel, err := fpsRepo.Get(flu.ProjectId)
		if utilities.IsValidError(err) {
			plog.Error("DB Error:", err)
			return errors.New("No Project Configuration found for FluProject:" + flu.ProjectId.String())
		}

		// reconsider
		maxFluCount := getMaxFluCount(fpsModel)
		postbackUrl := fpsModel.PostBackUrl
		//TODO Handle invalid url
		queryFrequency := getQueryFrequency(fpsModel)
		value = config{flu.ProjectId, fpsModel, maxFluCount, postbackUrl, queryFrequency}
		activeProjects[flu.ProjectId] = value
	}
	return nil
}
