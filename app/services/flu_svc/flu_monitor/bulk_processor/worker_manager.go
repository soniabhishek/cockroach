package bulk_processor

import (
	"time"
)

type WorkerManager struct {
	// A pool of workers channels that are registered with the dispatcher
	jobChannel chan jobChannel
	throttler  *time.Ticker

	jobChan jobChannel
	MaxJps int
}

func (wm *WorkerManager) Run() {
	go func() {
		for {
			<-wm.throttler.C

			job := <-wm.jobChan

			jobChannel := <-wm.jobChannel

			// a job request has been received
			go func(job1 Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle

				// dispatch the job to the worker job channel
				jobChannel <- job1
			}(job)
		}
	}()
}

func (wm *WorkerManager) PushJob(j Job) {
	wm.jobChan <- j
}

//jps (Jobs per second) - Worker manager will throttle job execution according if it crosses maxJps
func NewWorkerManager(maxJps int) WorkerManager {

	throttler := time.NewTicker(time.Duration(int(1000 / maxJps)) * time.Millisecond)
	return WorkerManager{jobChannel: make(chan jobChannel), throttler:throttler, jobChan:make(jobChannel), MaxJps:maxJps}
}
