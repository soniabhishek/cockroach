package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/utilities"
	"github.com/crowdflux/angel/app/config"
	"time"
)

var defaultFluThresholdCount = utilities.GetInt(config.DEFAULT_FLU_THRESHOLD_COUNT.Get())
var defaultRetryCount = utilities.GetInt(config.FLU_RETRY_THRESHOLD.Get())
var defaultRetryTimePeriod = time.Duration(utilities.GetInt(config.RETRY_TIME_PERIOD.Get())) * time.Millisecond
var fluThresholdDuration = int64(utilities.GetInt(config.FLU_THRESHOLD_DURATION.Get()))

func getQueryFrequency(fpsModel models.ProjectConfiguration) int {
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

func getRetryCount(fpsModel models.ProjectConfiguration) int {
	val := fpsModel.Options[RETRY_COUNT]
	if val == nil {
		return defaultRetryCount
	}
	retryCount := utilities.GetInt(val.(string))
	if retryCount == 0 {
		retryCount = defaultRetryCount
	}
	return retryCount
}

func getRetryPeriod(fpsModel models.ProjectConfiguration) time.Duration {
	val := fpsModel.Options[RETRY_COUNT]
	if val == nil {
		return defaultRetryTimePeriod
	}
	retryPeriod := time.Duration(utilities.GetInt(val.(string)))*time.Millisecond
	if retryPeriod == 0 {
		retryPeriod = defaultRetryTimePeriod
	}
	return retryPeriod
}
