package plog

import (
	// using library for logging
	"fmt"
	_ "github.com/mgutz/logxi/v1"
)

func Info(tag string, args ...interface{}) {
	if levelInfo <= plogLevel {
		fmt.Println(tag)
		fmt.Println(args)
	}
}
