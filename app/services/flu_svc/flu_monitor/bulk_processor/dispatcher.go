package bulk_processor

type Dispatcher struct {
	jobManagers []*JobManager
	workerPool  chan jobChannel
	maxWorkers  int
}

func NewDispatcher(maxWorkers int) *Dispatcher {

	return &Dispatcher{
		maxWorkers:    maxWorkers,
		workerPool:    make(chan jobChannel, maxWorkers),
	}
}

func (d *Dispatcher) AddJobManager(jm *JobManager) {
	d.jobManagers = append(d.jobManagers, jm)
}

func (d *Dispatcher) Start() (started bool) {

	d.startCheck()

	runJobManagers(d.jobManagers)

	d.startWorkers(d.maxWorkers)

	go d.dispatch()

	return true
}

func (d *Dispatcher) dispatch() {

	var jobChan jobChannel

	for {
		for _, jm := range d.jobManagers {

			// Wait for a worker's jobChannel from WorkerPool
			// if not available
			if jobChan == nil {
				jobChan = <-d.workerPool
			}

			// Pass that worker's jobChannel to any jobManager which is
			// ready to receive
			select {
			case jm.workerPool <- jobChan:
				jobChan = nil
			default:
			}
		}
	}
}

func runJobManagers(jobManagers []*JobManager) {

	for _, jm := range jobManagers {
		jm.Run()
	}
}

func (d *Dispatcher) startWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		w := newWorker(d.workerPool)
		w.Start()
	}
}

func (d *Dispatcher) startCheck() {
	if d.maxWorkers <= 0{
		panic("Max worker configured <= 0")
	}

	if len(d.jobManagers) == 0 {
		panic("No JobManagers added")
	}
}
