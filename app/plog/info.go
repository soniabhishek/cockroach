package plog

import (
	// using library for logging
	"fmt"
	_ "github.com/mgutz/logxi/v1"
	"runtime"
	"time"
)

func Info(tag string, args ...interface{}) {
	if levelInfo <= plogLevel {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(time.Now().Format(logFormat), fn, line, tag, args)
	}
}
