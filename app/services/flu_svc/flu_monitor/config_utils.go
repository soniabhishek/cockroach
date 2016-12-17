package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/utilities"
)

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
