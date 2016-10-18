package plog_logger

import "fmt"

type ConsoleLogger struct {
}

func (ConsoleLogger) Write(message ...interface{}){
	fmt.Println(message)
}
