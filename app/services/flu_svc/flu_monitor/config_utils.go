package flu_monitor

import (
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/utilities"
	"time"
)

var defaultFluThresholdCount = utilities.GetInt(config.DEFAULT_FLU_THRESHOLD_COUNT.Get())
var defaultRetryCount = utilities.GetInt(config.FLU_RETRY_THRESHOLD.Get())
var defaultClientQps = utilities.GetInt(config.DEFAULT_CLIENT_QPS.Get())
var defaultRetryTimePeriod = time.Duration(utilities.GetInt(config.RETRY_TIME_PERIOD.Get())) * time.Millisecond
var fluThresholdDuration = int64(utilities.GetInt(config.FLU_THRESHOLD_DURATION.Get()))

func getQueryFrequency(fpsModel models.ProjectConfiguration) int {
	val := fpsModel.Options[CLIENT_QPS]
	if val == nil {
		return defaultClientQps
	}
	queryFrequency := utilities.GetInt(val.(string))
	if queryFrequency == 0 {
		queryFrequency = defaultClientQps

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
