package plog

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
)

type CustomLogger struct {
	tag    string
	level  levelType
	logger *logrus.Logger
}

func NewLogger(name string, levelString string, typeString string) *CustomLogger {

	level, ok := levels[levelString]
	if !ok {
		panic("Invalid level string: " + levelString)
	}

	logger := logrus.New()
	setLogger(logger, path+"/"+name, typeString)
	return &CustomLogger{name, level, logger}
}

func (c *CustomLogger) Fatal(err error, args ...interface{}) {
	if levelFatal >= c.level {
		ErrorMail(c.tag, err, args)
		c.logger.WithFields(logrus.Fields{
			"error": err,
			"args":  fmt.Sprintf("%+v , ", args),
		}).Fatalln(c.tag)
	}
}

func (c *CustomLogger) Error(err error, args ...interface{}) {
	if levelError >= c.level {
		c.logger.WithFields(logrus.Fields{
			"error": err,
			"args":  fmt.Sprintf("%+v", args),
		}).Errorln(c.tag)
		ErrorMail(c.tag, err, args)
	}
}

func (c *CustomLogger) Warn(args ...interface{}) {
	if levelWarn >= c.level {
		c.logger.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v", args),
		}).Warnln(c.tag)
	}
}

func (c *CustomLogger) Info(args ...interface{}) {
	if levelInfo >= c.level {
		c.logger.WithFields(logrus.Fields{
			"args": fmt.Sprintf("%+v", args),
		}).Infoln(c.tag)
	}
}

func (c *CustomLogger) Debug(args ...interface{}) {

	if levelDebug >= c.level {
		_, fn, line, _ := runtime.Caller(1)
		c.logger.WithFields(logrus.Fields{
			"file": fn,
			"line": line,
			"args": fmt.Sprintf("%+v", args),
		}).Debugln(c.tag)
	}
}

func (c *CustomLogger) Trace(args ...interface{}) {

	if levelTrace >= c.level {
		_, fn, line, _ := runtime.Caller(1)
		c.logger.WithFields(logrus.Fields{
			"level": "trace",
			"file":  fn,
			"line":  line,
			"args":  fmt.Sprintf("%+v", args),
		}).Debugln(c.tag)
	}
}
