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

	//mutex.Lock // needed?
	//defer mutex.Unlock()
	for _, v := range activeProjects {
		var actualCount int
		if availabiltyPerClient>v.queryFrequency{
			actualCount = v.queryFrequency
		}else{
			actualCount = availabiltyPerClient
		}
		//make request in 1/actual_count time intervals
		availableQps -= actualCount
		rate := time.Second / time.Duration(actualCount)

		throttle := time.Tick(rate)
		for {
			<-throttle
			// spawn something instead?
			// spawn requests!! don't call them!

			//make the request in another pool? use channel!! send to channel!!
			// buffered channel?
			// go spawnRequest(v)

			// push to job queue
			// async create and then call the request object
			go makeRequest(v)
			// retry logic?
		}
	}

	// what does the channel do? accept from request pool make request
	// use a 'then' function to make repeated requests for retrying?

}
