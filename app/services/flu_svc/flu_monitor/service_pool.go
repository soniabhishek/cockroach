package flu_monitor

import "time"

func (fm *FluMonitor) servicePoolStart() error {
	if fm.PoolIsRunning {
		return
	}
	fm.PoolIsRunning = true

	rate := time.Second

	throttle := time.Tick(rate)
	for {
		<-throttle
		distributor() //call method to distribute every second
	}

	// make distributor and caller decoupled
}
