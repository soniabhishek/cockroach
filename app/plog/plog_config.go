package plog

import (
	"strings"
        "github.com/crowdflux/angel/app/plog/logger"
	"github.com/crowdflux/angel/app/config"
)

type levelType uint

const (
	levelNone levelType = iota + 1
	levelTrace
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
)

var levels = map[string]levelType{
	"NONE"  : levelNone,
	"TRACE" : levelTrace,
	"DEBUG" : levelDebug,
	"INFO"  : levelInfo,
	"WARN"  : levelWarn,
	"ERROR" : levelError,
	"FATAL" : levelFatal,
}

const (
	STR_TYPE_CONSOLE = "CONSOLE"
	STR_TYPE_FILE = "FILE"
)

var plogLevel levelType
var log plog_logger.ILogger

func init() {

	log = GetLogger()

	logLevelStr := strings.ToUpper(config.PLOG_LEVEL.Get())

	plogLevel= levels[logLevelStr]

        if plogLevel==0 {
		plogLevel = GetLevelFromEnvironment()
	}
}

func GetLogger()  plog_logger.ILogger{
	logTypeStr := strings.ToUpper(config.PLOG_TYPE.Get())
	switch logTypeStr {

	//TODO: Add file logger
	case STR_TYPE_CONSOLE:
		return plog_logger.ConsoleLogger{}
	default: return plog_logger.ConsoleLogger{}
	}
}

func GetLevelFromEnvironment() levelType {
	if config.IsDevelopment() {
		return levelTrace
	} else if config.IsStaging() {
		return levelInfo
	} else if config.IsProduction() {
		return levelError // Or we can put levelNone
	}

	return levelNone
}

var logFormat string = "2006-01-02 15:04:05.000"




