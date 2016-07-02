package plog

import (
	"fmt"
	"runtime"
)

func Trace(tag string, args ...interface{}) {

	if levelTrace <= plogLevel {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fn, line, tag, args)
	}
}
func IsTraceEnabled() bool {

	return levelTrace <= plogLevel
}
