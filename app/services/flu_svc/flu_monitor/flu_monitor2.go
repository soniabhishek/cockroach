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
	"math"
	"sync"
	"time"
)

type config struct {
	projectId      uuid.UUID
	maxFluCount    int
	postBackUrl    string
	queryFrequency int
	queue          feed_line.Fl
}

type FluMonitor struct {
	PoolIsRunning bool
}

type fluOutputStruct struct {
	ID          uuid.UUID   `json:"flu_id"`
	ReferenceId string      `json:"reference_id"`
	Tag         string      `json:"tag"`
	Status      string      `json:"status"`
	Result      interface{} `json:"result"`
}

var retryCount = make(map[uuid.UUID]int)
var mutex = &sync.RWMutex{}
var dbLogger = feed_line_repo.NewLogger()

var retryTimePeriod = time.Duration(utilities.GetInt(config.RETRY_TIME_PERIOD.Get())) * time.Millisecond
var fluThresholdDuration = int64(utilities.GetInt(config.FLU_THRESHOLD_DURATION.Get()))
var monitorTimePeriod = time.Duration(utilities.GetInt(config.MONITOR_TIME_PERIOD.Get())) * time.Millisecond
var retryThreshold = utilities.GetInt(config.FLU_RETRY_THRESHOLD.Get())

var projectConfig = make(map[uuid.UUID]config) // Hash map to store config
var queues = make(map[uuid.UUID]feed_line.Fl)  // Hash map to store queues

var defaultFluThresholdCount = utilities.GetInt(config.DEFAULT_FLU_THRESHOLD_COUNT.Get())
var totalQps = utilities.GetInt(config.TOTAL_QPS.Get())
var availableQps = totalQps

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	clientQ := queues[flu.ProjectId]
	if clientQ == nil {
		clientQ := feed_line.New(flu.ProjectId.String())
		queues[flu.ProjectId] = clientQ
	}
	clientQ.Push(flu)

	saveProjectConfig(flu)

	fm.servicePoolStart()
	return
}

func saveProjectConfig(flu models.FeedLineUnit) {
	value, valuePresent := projectConfig[flu.ProjectId]
	if valuePresent == false {
		fpsRepo := project_configuration_repo.New()
		fpsModel, err := fpsRepo.Get(flu.ProjectId)
		if utilities.IsValidError(err) {
			plog.Error("DB Error:", err)
			return errors.New("No Project Configuration found for FluProject:" + flu.ProjectId.String())
		}
		maxFluCount := getMaxFluCount(fpsModel)
		postbackUrl := fpsModel.PostBackUrl
		//TODO Handle invalid url
		queryFrequency := getQueryFrequency(fpsModel)
		value = config{flu.ProjectId, maxFluCount, postbackUrl, queryFrequency}
		projectConfig[flu.ProjectId] = value
	}
}

func getQueryFrequency(fpsModel models.ProjectConfiguration) interface{} {
	val := fpsModel.Options[QUERY_FREQUENCY]
	if val == nil {
		//TODO change later. take from config
		return 5
	}
	queryFrequency := utilities.GetInt(val.(string))
	if queryFrequency == 0 {
		//TODO change later. take from config
		queryFrequency = 5
	}
	return queryFrequency
}

func getMaxFluCount(fpsModel models.ProjectConfiguration) int {
	val := fpsModel.Options[MAX_FLU_COUNT]
	if val == nil {
		return defaultFluThresholdCount
	}
	maxFluCount := utilities.GetInt(val.(string))
	if maxFluCount == 0 {
		maxFluCount = defaultFluThresholdCount
	}
	return maxFluCount
}

func (fm *FluMonitor) servicePoolStart() error {
	if fm.PoolIsRunning {
		return
	}
	fm.PoolIsRunning = true

	rate := time.Second

	throttle := time.Tick(rate)
	for {
		<-throttle
		distributor() //call method to distribute every second
	}
}

func distributor() {

	// get clients count
	clientCount := len(projectConfig)
	// divide our capacity/number of clients = somenum
	availabiltyPerClient := availableQps / clientCount

	mutex.Lock // needed?
	defer mutex.Unlock()
	for k, v := range projectConfig {
		actualCount := math.Min(availabiltyPerClient, v.queryFrequency)
		//make request in 1/actual_count time intervals
		availableQps -= actualCount
		rate := time.Second / actualCount

		throttle := time.Tick(rate)
		for {
			<-throttle
			go func() {
				makeRequest(k, v)
			}()
			// retry logic?
		}
	}

}
func makeRequest(projectId uuid.UUID, projectConfig config) error {
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
		if flu == nil {
			delete(projectConfig, projectId)
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

	return
}

/*func getFluOutputObj(projectConfig config) (fluOutputObj []fluOutputStruct) {

	return
}*/
