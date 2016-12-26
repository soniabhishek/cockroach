package flu_monitor

import "net/http"

type DoFunc func()

type Job struct {
	Request http.Request
	Do1 DoFunc
	RetryCount int
	RetryInterval
}


// A buffered channel that we can send work requests on.
var JobQueue chan Job = make(chan Job)