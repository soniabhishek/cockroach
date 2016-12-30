package flu_monitor

import (
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

var defaultClientQps = services.AtoiOrPanic(config.DEFAULT_CLIENT_QPS.Get())
var defaultMaxFluCount = 1
var retryCount = services.AtoiOrPanic(config.FLU_RETRY_THRESHOLD.Get())
var retryTimePeriod = time.Duration(services.AtoiOrPanic(config.RETRY_TIME_PERIOD.Get())) * time.Millisecond

func getHmacHeader(pcModel models.ProjectConfiguration) (hmacHeader string) {

	val := pcModel.Options[HMAC_HEADER_KEY]
	if val == nil {
		return
	}
	valString, ok := val.(string)
	if !ok {
		plog.Error("Flu monitor", errors.New("error parsing hmac_header from project_configuration. Not string. Using default."))
		return
	}
	return valString
}

func getHmacKey(pcModel models.ProjectConfiguration) (hmacKey string) {

	val := pcModel.Options[HMAC_KEY]
	if val == nil {
		return
	}
	valString, ok := val.(string)
	if !ok {
		plog.Error("Flu monitor", errors.New("error parsing hmac_key from project_configuration. Not string. Using default."))
		return
	}
	return valString
}

func isHmacEnabled(pcModel models.ProjectConfiguration) bool {

	val := pcModel.Options[IS_HMAC_ENABLED]
	if val == nil {
		return false
	}
	valBool, ok := val.(bool)
	if !ok {
		plog.Error("Flu monitor", errors.New("error parsing is_hmac_enabled from project_configuration. Not string. Using default."))
		return false
	}
	return valBool
}

func getQueryFrequency(pcModel models.ProjectConfiguration) int {
	val := pcModel.Options[CLIENT_QPS]
	if val == nil {
		return defaultClientQps
	}

	valString, ok := val.(string)
	if !ok {
		plog.Error("Flu monitor", errors.New("error parsing client_qps from project_configuration. Not string. Using default."))
		return defaultClientQps
	}

	queryFrequency, err := strconv.Atoi(valString)
	if err != nil {
		plog.Error("Flu monitor", errors.New("error parsing client_qps from project_configuration. Invalid string. Using default."))
		queryFrequency = defaultClientQps

	}
	return queryFrequency
}

func getMaxFluCount(pcModel models.ProjectConfiguration) int {
	val := pcModel.Options[MAX_FLU_COUNT]
	if val == nil {
		return defaultMaxFluCount
	}

	valString, ok := val.(string)
	if !ok {
		plog.Error("Flu monitor", errors.New("error parsing max_flu_count from project_configuration. Not string. Using default."))
		return defaultMaxFluCount
	}

	queryFrequency, err := strconv.Atoi(valString)
	if err != nil {
		plog.Error("Flu monitor", errors.New("error parsing max_flu_count from project_configuration. Invalid string. Using default."))
		queryFrequency = defaultMaxFluCount

	}
	return queryFrequency
}
