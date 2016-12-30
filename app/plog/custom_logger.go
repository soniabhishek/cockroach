package plog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
)

type CustomLogger struct {
	tag   string
	level levelType
}

func NewLogger(name string, levelString string) *CustomLogger {

	level, ok := levels[levelString]
	if !ok {
		panic("Invalid level string: " + levelString)
	}

	return &CustomLogger{name, level}
}

func (c *CustomLogger) Fatal(err error, args ...interface{}) {
	if levelFatal >= c.level {
		ErrorMail(c.tag, err, args)
		logr.WithFields(logrus.Fields{
			"error": err,
			"args":  fmt.Sprintf("%+v , ", args),
		}).Fatalln(c.tag)
	}
}

func (c *CustomLogger) Error(err error, args ...interface{}) {
	if levelError >= c.level {
		logr.WithFields(logrus.Fields{
			"error": err,
			"args":  fmt.Sprintf("%+v", args),
		}).Errorln(c.tag)
		ErrorMail(c.tag, err, args)
	}
}

func (c *CustomLogger) Warn(args ...interface{}) {
	if levelWarn >= c.level {
		logr.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v", args),
		}).Warnln(c.tag)
	}
}

func (c *CustomLogger) Info(args ...interface{}) {
	if levelInfo >= c.level {
		logr.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v", args),
		}).Infoln(c.tag)
	}
}

func (c *CustomLogger) Debug(args ...interface{}) {

	if levelDebug >= c.level {
		_, fn, line, _ := runtime.Caller(1)
		logr.WithFields(logrus.Fields{
			"file": fn,
			"line": line,
			"args": fmt.Sprintf("%+v", args),
		}).Debugln(c.tag)
	}
}

func (c *CustomLogger) Trace(args ...interface{}) {

	if levelTrace >= c.level {
		_, fn, line, _ := runtime.Caller(1)
		logr.WithFields(logrus.Fields{
			"level": "trace",
			"file":  fn,
			"line":  line,
			"args":  fmt.Sprintf("%+v", args),
		}).Debugln(c.tag)
	}
}
