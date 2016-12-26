package bulk_processor

import (
	"github.com/crowdflux/angel/app/plog"
	"fmt"
)

// Worker represents the worker that executes the job
type Worker struct {
	//TODO use JobChan
	WorkerPool chan jobChannel
	JobChannel chan Job
	quit       chan bool
}

func (w *Worker) Start() {

	fmt.Println("worker started")

	go func() {
		for {

			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:

				plog.Trace("Bulk Processor", "Worker", "Starting Job")
				// we have received a work request.
				job.Do()

				plog.Trace("Bulk Processor", "Worker", "Finished Job")

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func newWorker(workerPool chan jobChannel) *Worker {
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(jobChannel),
		quit:       make(chan bool)}
}
