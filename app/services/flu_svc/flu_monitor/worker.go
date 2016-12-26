package flu_monitor

// Worker represents the worker that executes the job
type Worker struct {

	//TODO use JobCham
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func (w Worker) Start() {
	go func() {
		for {

			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				job.Do()
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

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}