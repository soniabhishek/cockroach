package bulk_processor

type Dispatcher struct {
	workerManagers []WorkerManager
	workerPool     chan jobChannel
	maxWorkers     int
}

func NewDispatcher(maxWorkers int) *Dispatcher {

	return &Dispatcher{
		maxWorkers:    maxWorkers,
		workerPool:    make(chan jobChannel, maxWorkers),
	}
}

func (d *Dispatcher) AddWorkerManager(m WorkerManager) {
	d.workerManagers = append(d.workerManagers, m)
}

func (d *Dispatcher) Start() (started bool) {

	runManagers(d.workerManagers)
	d.startWorkers(d.maxWorkers)

	go d.dispatch()

	return true
}

func (d *Dispatcher) dispatch() {

	var tempJobChan jobChannel

	for {
		if tempJobChan == nil {
			tempJobChan = <-d.workerPool
		}

		loop:
		for _, wm := range d.workerManagers {
			select {
			case wm.jobChannel <- tempJobChan:
				tempJobChan = nil
				break loop
			default:
			}
		}
	}
}

func runManagers(workerManagers []WorkerManager) {

	if workerManagers == nil {
		panic("no workerManagers found")
	}

	for _, wm := range workerManagers {
		wm.Run()
	}
}
func (d *Dispatcher) startWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		w := newWorker(d.workerPool)
		w.Start()
	}
}
