package plog

import (
	"fmt"
)

func Trace(tag string, args ...interface{}) {

	if levelTrace <= plogLevel {
		fmt.Println(tag)
		fmt.Println(args)
	}
}
