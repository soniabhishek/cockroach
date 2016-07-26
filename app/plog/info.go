package plog

import (
	// using library for logging
	"fmt"
	"runtime"
	"time"

	_ "github.com/mgutz/logxi/v1"
)

func Info(tag string, args ...interface{}) {
	if levelInfo <= plogLevel {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(time.Now().Format(logFormat), fn, line, tag, args)
	}
}
