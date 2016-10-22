package plog

import (
	"strings"
        "github.com/crowdflux/angel/app/config"
	"github.com/Sirupsen/logrus"
	"os"
	"time"
	"fmt"
	"path/filepath"
	"github.com/jasonlvhit/gocron"
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
var logr = logrus.New()
var path,_= filepath.Abs("./app_logs")

func init() {

	SetLogger()

	logLevelStr := strings.ToUpper(config.PLOG_LEVEL.Get())

	plogLevel= levels[logLevelStr]

        if plogLevel==0 {
		plogLevel = GetLevelFromEnvironment()
	}
}

func SetLogger() {

	logr.Formatter = &logrus.JSONFormatter{}

	logr.Level = logrus.DebugLevel

	_ = os.Mkdir(path,os.ModePerm)

	SetLogOutput()
	log_file_scheduler := gocron.NewScheduler()
	log_file_scheduler.Every(1).Day().At("00.00").Do(SetLogOutput)
	log_file_scheduler.Start()
}
func SetLogOutput() {

		logTypeStr := strings.ToUpper(config.PLOG_TYPE.Get())

		switch logTypeStr {
		case STR_TYPE_CONSOLE:
			logr.Out = os.Stdout
		case STR_TYPE_FILE:
			logr.Out = GetFileIO(path + "/" + GetFileName())
		case STR_TYPE_ERROR:
			logr.Out = os.Stderr

		default: logr.Out = os.Stdout
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

func GetFileName() string{
	y,m,d := time.Now().Date()
	dateString:= fmt.Sprintf("%d_%d_%d", y,m,d)
	return "log_"+dateString+".txt"
}

func GetFileIO(path string) *os.File{
	CreateFile(path)
	ret,_ := os.OpenFile(path,os.O_RDWR|os.O_APPEND,0660)
	return ret
}

func CreateFile(path string) {

	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, _ = os.Create(path)
		defer file.Close()
	}
}





