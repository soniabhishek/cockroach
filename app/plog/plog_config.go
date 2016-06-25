package plog

import (
	"gitlab.com/playment-main/angel/app/config"
	"strings"
)

const (
	levelNone  = -1
	levelError = 1
	levelInfo  = 3
	levelTrace = 5
)

const (
	STR_NONE  = "NONE"
	STR_ERROR = "ERROR"
	STR_INFO  = "INFO"
	STR_TRACE = "TRACE"
)

var plogLevel int

func init() {
	logLevelStr := strings.ToUpper(config.PLOG_LEVEL.Get())
	switch logLevelStr {

	case STR_NONE:
		plogLevel = levelNone
		break
	case STR_ERROR:
		plogLevel = levelError
		break
	case STR_INFO:
		plogLevel = levelInfo
		break
	case STR_TRACE:
		plogLevel = levelTrace
		break

	default:
		if config.IsDevelopment() {
			plogLevel = levelTrace
		} else if config.IsStaging() {
			plogLevel = levelInfo
		} else if config.IsProduction() {
			plogLevel = levelError // Or we can put levelNone
		} else {
			plogLevel = levelNone
		}
	}
}
