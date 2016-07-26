package crowdsourcing_step

import (
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/work_flow_svc/counter"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type crowdSourcingStep struct {
	step.Step
	fluRepo   feed_line_repo.IFluRepo
	fluClient fluPusher
}

// Rename the interface later
type fluPusher interface {
	PushFLU(models.FeedLineUnit) (bool, error)
}

func (c *crowdSourcingStep) processFlu(flu feed_line.FLU) {

	c.AddToBuffer(flu)

	ok, err := c.fluClient.PushFLU(flu.FeedLineUnit)

	if !ok {
		c.Detain(flu, err, c.fluRepo)
	}
}

func (c *crowdSourcingStep) finishFlu(flu feed_line.FLU) bool {

	err := c.RemoveFromBuffer(flu)
	if err != nil {
		return false
	}
	counter.Print(flu, "crowdsourcing")
	c.OutQ <- flu
	return true
}

func (c *crowdSourcingStep) start() {
	go func() {
		for {
			select {
			case flu := <-c.InQ:
				c.processFlu(flu)
			}
		}
	}()
}

func (c *crowdSourcingStep) Connect(routerIn *feed_line.Fl) (routerOut *feed_line.Fl) {

	// Send output of this step to the router's input
	// for next rerouting
	c.OutQ = *routerIn

	c.start()

	// Return the input channel of this step
	// so that router can push flu to it
	return &c.InQ
}
