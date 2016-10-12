package router_svc

import (
	"errors"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
)

func start(sr *stepRouter) {

	// Start listening for incoming flus from InQ channel
	// in another goroutine & route it to its exact step
	go func() {

		for flu := range sr.InQ.Receiver() {

			plog.Info("Router in", flu.ID)
			// There is a question that adding to the
			// buffer should be inside or outside
			// the below go routine.
			//
			// Current implementation causes InQ to block if channel is full.
			// If put inside the below go routine InQ will become non blocking
			// but there is a chance of large number of go routines at a time
			sr.buffer <- 1

			// need to pass the flu as the function parameter as
			// the reference gets overridden in the next loop
			go func(flu feed_line.FLU) {

				defer func() {

					if r := recover(); r != nil {
						plog.Error("Router", errors.New("Panic occured in router"), r)
						(*sr.routeTable[step_type.Manual]).Push(flu)
						plog.Info("Router", "Sent to manual after panic", flu)
					}

					// Free the buffer
					<-sr.buffer
				}()

				// Add workers here
				// Right now its like just one sync worker
				// i.e. if the below method takes 1 second & buffer size is 10
				// then the max speed of router processing will be 1 flu * buffer/second = 10 flu/second
				r := sr.getRoute(&flu)
				plog.Info("router", "got some route")
				(*r).Push(flu)
				plog.Info("router out", "sent somewhere")
			}(flu)
		}
	}()
}
