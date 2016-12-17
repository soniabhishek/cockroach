package flu_monitor

import (
	"math"
	"time"
)

func distributor() {

	// get clients count
	clientCount := len(activeProjects)
	// divide our capacity/number of clients = availabilityPerClient
	availabiltyPerClient := availableQps / clientCount

	mutex.Lock // needed?
	defer mutex.Unlock()
	for _, v := range activeProjects {
		actualCount := math.Min(availabiltyPerClient, v.queryFrequency)
		//make request in 1/actual_count time intervals
		availableQps -= actualCount
		rate := time.Second / actualCount

		throttle := time.Tick(rate)
		for {
			<-throttle
			go makeRequest(v)
			// retry logic?
		}
	}

}
