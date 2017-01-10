package plog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/crowdflux/angel/app/config"
	"github.com/getsentry/raven-go"
	"runtime"
)

func Fatal(tag string, err error, args ...interface{}) {
	if levelFatal >= plogLevel {
		ErrorMail(tag, err, args)
		logr.WithFields(logrus.Fields{
			"error": err,
			"args":  fmt.Sprintf("%+v , ", args),
		}).Fatalln(tag)
	}
}

// we can pass plog.Message here in place if raven.Interface
func Error(tag string, err error, args ...message) {
	if levelError >= plogLevel {
		if !config.IsProduction() {
			sentryItems := map[string]string{"tag": tag}
			for _, arg := range args {
				key := string(arg.Tag.Type)
				value := fmt.Sprintf("%+v", arg.Params)

				val, ok := sentryItems[key]
				if !ok {
					sentryItems[key] = value

				} else {
					sentryItems[key] = val + " ; " + value
				}
			}
			raven.DefaultClient.CaptureError(err, sentryItems)
		}
		logr.WithFields(logrus.Fields{
			"error": err,
			"args":  fmt.Sprintf("%+v", args),
		}).Errorln(tag)
		ErrorMail(tag, err, args)
	}
}

func Warn(tag string, args ...interface{}) {
	if levelWarn >= plogLevel {
		logr.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v", args),
		}).Warnln(tag)
	}
}

func Info(tag string, args ...interface{}) {
	if levelInfo >= plogLevel {
		logr.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v", args),
		}).Infoln(tag)
	}
}

func Debug(tag string, args ...interface{}) {

	if levelDebug >= plogLevel {
		_, fn, line, _ := runtime.Caller(1)
		logr.WithFields(logrus.Fields{
			"file": fn,
			"line": line,
			"args": fmt.Sprintf("%+v", args),
		}).Debugln(tag)
	}
}

func Trace(tag string, args ...interface{}) {

	if levelTrace >= plogLevel {
		_, fn, line, _ := runtime.Caller(1)
		logr.WithFields(logrus.Fields{
			"level": "trace",
			"file":  fn,
			"line":  line,
			"args":  fmt.Sprintf("%+v", args),
		}).Debugln(tag)
	}
}
