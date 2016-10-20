package plog

import (
	"runtime"
	"os"
)

func Fatal(tag string, err error, args ...interface{}) {
	if levelFatal <= plogLevel {
		logr.Fatalf(tag, err, args)
		ErrorMail(tag, err, args)
		os.Exit(1)
	}
}

func Error(tag string, err error, args ...interface{}) {
	if levelError <= plogLevel {
		logr.Errorf(tag, err, args)
		ErrorMail(tag, err, args)
	}
}

func Warn(tag string, args ...interface{}) {
	if levelWarn <= plogLevel {
		logr.Warnf(tag, args)
	}
}

func Info(tag string, args ...interface{}) {
	if levelInfo <= plogLevel {
		logr.Infof(tag,args)
	}
}


func Debug(tag string, args ...interface{}) {

	if levelDebug <= plogLevel{
		_, fn, line, _ := runtime.Caller(1)
		logr.Debugf(tag,fn,line, args)
	}
}

func Trace(tag string, args ...interface{}) {

	if IsTraceEnabled() {
		_, fn, line, _ := runtime.Caller(1)

		logr.Debugf(" TRACE: "+tag,fn,line, args)

	}
}

func IsTraceEnabled() bool {
	return levelTrace <= plogLevel
}
