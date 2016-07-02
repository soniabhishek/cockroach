package plog

import (
	"fmt"
	"runtime"
	"time"
)

func Trace(tag string, args ...interface{}) {

	if levelTrace <= plogLevel {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(time.Now().Format(logFormat), fn, line, tag, args)
	}
}
func IsTraceEnabled() bool {
	return levelTrace <= plogLevel
}
