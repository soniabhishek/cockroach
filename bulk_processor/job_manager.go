package bulk_processor

import (
	"time"
)

type JobManager struct {
	// A pool of workers channels that are registered with the dispatcher
	allocatedWorker chan jobChannel

	// Used to throttle
	throttler <-chan time.Time

	// To communicate from PushJob() To Run()
	jobChan jobChannel

	MaxJps int
	name   string
}

func (jm *JobManager) Run() {
	go func() {
		for {
			// Throttle the loop according to input Jps
			<-jm.throttler

			// Wait for a job from PushJob()
			job := <-jm.jobChan

			// Get a worker (type JobChannel) from Worker Pool
			jobChannel := <-jm.allocatedWorker

			// Push the job to job channel
			jobChannel <- job

		}
	}()
}

// A blocking call to JobManager
func (jm *JobManager) PushJob(j Job) {
	jm.jobChan <- j
}

//jps (Jobs per second) - JobManager will throttle job execution according if it crosses maxJps
func NewJobManager(maxJps int, name string) *JobManager {

	throttler := time.Tick(time.Duration(int(1000/maxJps)) * time.Millisecond)

	return &JobManager{
		// Zero size WorkerPool
		allocatedWorker: make(chan jobChannel),
		// To throttle
		throttler: throttler,
		// For communicating from PushJob() To Run()
		jobChan: make(jobChannel),
		MaxJps:  maxJps,
		name:    name,
	}
}
