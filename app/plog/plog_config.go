package plog

import (
	"strings"
        "github.com/crowdflux/angel/app/config"
	"github.com/Sirupsen/logrus"
	"os"
	"time"
	"fmt"
	"path/filepath"
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
	"NONE"  : levelNone,
	"TRACE" : levelTrace,
	"DEBUG" : levelDebug,
	"INFO"  : levelInfo,
	"WARN"  : levelWarn,
	"ERROR" : levelError,
	"FATAL" : levelFatal,
	"ALL"	: levelAll,
}

const (
	STR_TYPE_CONSOLE = "CONSOLE"
	STR_TYPE_FILE = "FILE"
	STR_TYPE_ERROR = "ERROR"
)

var plogLevel levelType
var logr *logrus.Logger

func init() {

	logr = GetLogger()

	logLevelStr := strings.ToUpper(config.PLOG_LEVEL.Get())

	plogLevel= levels[logLevelStr]

        if plogLevel==0 {
		plogLevel = GetLevelFromEnvironment()
	}
}

func GetLogger() *logrus.Logger{
	logTypeStr := strings.ToUpper(config.PLOG_TYPE.Get())

	var logr = logrus.New()

	logr.Level = logrus.DebugLevel

	path,_ := filepath.Abs("./app_logs")

	_ = os.Mkdir(path,os.ModePerm)

	switch logTypeStr {

	//TODO: Add file logger
	case STR_TYPE_CONSOLE:
		logr.Out = os.Stdout
	case STR_TYPE_FILE:
		logr.Out = getFileIO(path+"/"+getFileName())
	case STR_TYPE_ERROR:
		logr.Out = os.Stderr

	default: logr.Out = os.Stdout
	}
	return logr
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

func getFileName() string{
	y,m,d := time.Now().Date()
	dateString:= fmt.Sprintf("%d_%d_%d", y,m,d)
	return "log_"+dateString+".txt"
}

func getFileIO(path string) *os.File{
	createFile(path)
	ret,_ := os.OpenFile(path,os.O_RDWR|os.O_APPEND,0660)
	return ret
}

func createFile(path string) {

	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, _ = os.Create(path)
		defer file.Close()
	}
}





