package plog

import (
	"runtime"
	"os"
)

func Fatal(tag string, err error, args ...interface{}) {
	if levelFatal >= plogLevel {
		logr.Fatalln(tag, err, args)
		ErrorMail(tag, err, args)
		os.Exit(1)
	}
}

func Error(tag string, err error, args ...interface{}) {
	if levelError >= plogLevel {
		logr.Errorln(tag, err, args)
		ErrorMail(tag, err, args)
	}
}

func Warn(tag string, args ...interface{}) {
	if levelWarn >= plogLevel {
		logr.Warnln(tag, args)
	}
}

func Info(tag string, args ...interface{}) {
	if levelInfo >= plogLevel {
		logr.Infoln(tag,args)
	}
}


func Debug(tag string, args ...interface{}) {

	if levelDebug >= plogLevel{
		_, fn, line, _ := runtime.Caller(1)
		logr.Debugln(tag,fn,line, args)
	}
}

func Trace(tag string, args ...interface{}) {

	if IsTraceEnabled() {
		_, fn, line, _ := runtime.Caller(1)

		logr.Debugln(" TRACE: "+tag,fn,line, args)

	}
}

func IsTraceEnabled() bool {
	return levelTrace >= plogLevel
}
