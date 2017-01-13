package plog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/crowdflux/angel/app/config"
	"github.com/jasonlvhit/gocron"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	levelAll
)

var levels = map[string]levelType{
	"NONE":  levelNone,
	"TRACE": levelTrace,
	"DEBUG": levelDebug,
	"INFO":  levelInfo,
	"WARN":  levelWarn,
	"ERROR": levelError,
	"FATAL": levelFatal,
	"ALL":   levelAll,
}

const (
	STR_TYPE_CONSOLE = "CONSOLE"
	STR_TYPE_FILE    = "FILE"
	STR_TYPE_ERROR   = "ERROR"
)

var plogLevel levelType
var logr = logrus.New()
var path, _ = filepath.Abs(config.PLOG_LOCATION.Get())

func init() {

	logTypeStr := strings.ToUpper(config.PLOG_TYPE.Get())

	setLogger(logr, path, logTypeStr)

	logLevelStr := strings.ToUpper(config.PLOG_LEVEL.Get())

	plogLevel = levels[logLevelStr]

	if plogLevel == 0 {
		plogLevel = getLevelFromEnvironment()
	}
}

func setLogger(logger *logrus.Logger, logPath string, logType string) {

	logger.Formatter = &logrus.JSONFormatter{}

	logger.Level = logrus.DebugLevel

	_ = os.Mkdir(logPath, os.ModePerm)

	setLogOutput(logger, logType)
}
func setLogOutput(logger *logrus.Logger, logTypeStr string) {

	switch logTypeStr {
	case STR_TYPE_CONSOLE:
		logger.Out = os.Stdout
	case STR_TYPE_FILE:
		setFileIO(logger)
		log_file_scheduler := gocron.NewScheduler()
		log_file_scheduler.Every(1).Day().At("00.00").Do(setFileIO, logger)
		log_file_scheduler.Start()
	case STR_TYPE_ERROR:
		logger.Out = os.Stderr

	default:
		logger.Out = os.Stdout
	}
}

func getLevelFromEnvironment() levelType {
	if config.IsDevelopment() {
		return levelTrace
	} else if config.IsStaging() {
		return levelInfo
	} else if config.IsProduction() {
		return levelError // Or we can put levelNone
	}

	return levelNone
}

func getFileName() string {
	y, m, d := time.Now().Date()
	dateString := fmt.Sprintf("%d_%d_%d", y, m, d)
	return "log_" + dateString + ".txt"
}

func setFileIO(logger *logrus.Logger) {
	file_location := path + "/" + getFileName()
	createFile(file_location)
	logger.Out, _ = os.OpenFile(file_location, os.O_RDWR|os.O_APPEND, 0660)
}

func createFile(file_location string) {

	// detect if file exists
	var _, err = os.Stat(file_location)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, _ = os.Create(file_location)
		defer file.Close()
	}
}
