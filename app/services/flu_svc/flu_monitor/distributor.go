package flu_monitor

import (
	"time"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
)

func distributor() {


	// get clients count
	clientCount := len(activeProjects)
	// divide our capacity/number of clients = availabilityPerClient
	availabiltyPerClient := availableQps / clientCount

	//mutex.Lock // needed?
	//defer mutex.Unlock()
	for _, v := range activeProjects {
		var finalQps int
		if availabiltyPerClient > v.queryFrequency{
			finalQps = v.queryFrequency
		}else{
			finalQps = availabiltyPerClient
		}
		//make request in 1/actual_count time intervals
		availableQps -= finalQps
		rate := time.Second / time.Duration(finalQps)

		throttle := time.Tick(rate)
		for {
			<-throttle
			go makeRequest(v)
			// retry logic?
		}
	}

}
