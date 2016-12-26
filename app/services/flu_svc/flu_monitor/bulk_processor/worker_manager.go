package bulk_processor

import (
	"time"
	//"fmt"
)

type WorkerManager struct {
	// A pool of workers channels that are registered with the dispatcher
	jobChannel chan jobChannel
	throttler  <- chan time.Time

	jobChan jobChannel
	MaxJps  int
	name    string
}

func (wm *WorkerManager) Run() {
	go func() {
		for {
			// Throttle the loop according to input Jps
			<-wm.throttler

			// Get the job from PushJob()
			job := <-wm.jobChan

			// Get the job channel
			jobChannel := <-wm.jobChannel

			// Push the job to job channel
			jobChannel <- job
		}
	}()
}

// A blocking call to WorkerManager
func (wm *WorkerManager) PushJob(j Job) {
	wm.jobChan <- j
}

//jps (Jobs per second) - Worker manager will throttle job execution according if it crosses maxJps
func NewWorkerManager(maxJps int, name string) *WorkerManager {

	throttler := time.Tick(time.Duration(int(1000 / maxJps)) * time.Millisecond)

	return &WorkerManager{
		jobChannel: make(chan jobChannel),
		throttler:  throttler,
		jobChan:    make(jobChannel),
		MaxJps:     maxJps,
		name:       name,
	}
}
