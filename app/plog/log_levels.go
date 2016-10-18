package plog

import (
	"runtime"
	"time"
	"os"
)

func Fatal(tag string, err error, args ...interface{}) {
	if levelFatal <= plogLevel {
		log.Write(time.Now().Format(logFormat) + " FATAL: ")
		ErrorMail(tag, err, args)
		os.Exit(1)
	}
}

func Error(tag string, err error, args ...interface{}) {
	if levelError <= plogLevel {
		log.Write(time.Now().Format(logFormat) + " ERROR: ")
		ErrorMail(tag, err, args)
	}
}

func Warn(tag string, args ...interface{}) {
	if levelWarn <= plogLevel {
		log.Write(time.Now().Format(logFormat) + " WARN : "+ tag, args)
	}
}

func Info(tag string, args ...interface{}) {
	if levelInfo <= plogLevel {
		log.Write(time.Now().Format(logFormat) + " INFO : "+tag, args)
	}
}


func Debug(tag string, args ...interface{}) {

	if levelDebug <= plogLevel{
		_, fn, line, _ := runtime.Caller(1)
		log.Write(time.Now().Format(logFormat) + " DEBUG : ", fn, line, tag, args)
	}
}

func Trace(tag string, args ...interface{}) {

	if IsTraceEnabled() {
		_, fn, line, _ := runtime.Caller(1)
		log.Write(time.Now().Format(logFormat) + " TRACE : ", fn, line, tag, args)
	}
}

func IsTraceEnabled() bool {
	return levelTrace <= plogLevel
}

