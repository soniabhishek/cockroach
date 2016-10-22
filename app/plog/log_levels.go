package plog

import (
	"runtime"
	"github.com/Sirupsen/logrus"
	"fmt"
)

func Fatal(tag string, err error, args ...interface{}) {
	if levelFatal >= plogLevel {
		ErrorMail(tag, err, args)
		logr.WithFields(logrus.Fields{
			"error" : err,
			"args"  : fmt.Sprintf("%+v , ",args),
		}).Fatalln(tag)
	}
}

func Error(tag string, err error, args ...interface{}) {
	if levelError >= plogLevel {
		logr.WithFields(logrus.Fields{
			"error" : err,
			"args"  : fmt.Sprintf("%+v",args),
		}).Errorln(tag)
		ErrorMail(tag, err, args)
	}
}

func Warn(tag string, args ...interface{}) {
	if levelWarn >= plogLevel {
		logr.WithFields(logrus.Fields{
			"args" : fmt.Sprintf("%+v",args),
		}).Warnln(tag)
	}
}

func Info(tag string, args ...interface{}) {
	if levelInfo >= plogLevel {
		logr.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v",args),
		}).Infoln(tag)
	}
}


func Debug(tag string, args ...interface{}) {

	if levelDebug >= plogLevel{
		_, fn, line, _ := runtime.Caller(1)
		logr.WithFields(logrus.Fields{
			"file" : fn,
			"line" : line,
			"args" : fmt.Sprintf("%+v",args),
		}).Debugln(tag)
	}
}

func Trace(tag string, args ...interface{}) {

	if IsTraceEnabled() {
		_, fn, line, _ := runtime.Caller(1)
		logr.WithFields(logrus.Fields{
			"level": "trace",
			"file" : fn,
			"line" : line,
			"args" : fmt.Sprintf("%+v",args),
		}).Debugln(tag)
	}
}

func IsTraceEnabled() bool {
	return levelTrace >= plogLevel
}
