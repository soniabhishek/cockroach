package logger

import (
	"fmt"

	// using library for logging
	_ "github.com/mgutz/logxi/v1"
)

func ErrorNMail(tag string, message string, err error, args ...interface{}) {
	fmt.Println(tag, message, err, args)
}
