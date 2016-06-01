package plogger

import (
	"fmt"

	// using library for logging
	_ "github.com/mgutz/logxi/v1"
)

func ErrorNMail(tag string, err error, message string, args ...interface{}) {
	fmt.Println(tag, err, message, args)
}
