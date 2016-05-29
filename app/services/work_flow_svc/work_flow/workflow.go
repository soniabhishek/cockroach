package work_flow

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step_router"
)

type WorkFlow struct {
	InQ  feed_line.FL
	OutQ feed_line.FL
}

// Creates a new workFlow instance
// Making it private because making it public can be dangerous
func newStdWorkFlow() WorkFlow {

	//create new instance
	w := WorkFlow{
		InQ:  feed_line.New(),
		OutQ: feed_line.New(),
	}

	// Start Workflow Channel IO in another goroutine
	go func() {
		for {
			select {

			case flu := <-w.InQ:
				flu.Step = "workflow in"
				counter.Print(flu)
				step_router.StdStepRouter.InQ <- flu

			case flu := <-step_router.StdStepRouter.ProcessedFluQ:
				flu.Step = "workflow out"
				counter.Print(flu)
				w.OutQ <- flu
			}
		}
	}()

	return w
}

// Exposing a StdWorkFlow instance
var StdWorkFlow = newStdWorkFlow()

//var StdWorkFlow = newShortCircut()

func newShortCircut() WorkFlow {
	//create new instance
	w := WorkFlow{
		InQ:  feed_line.New(),
		OutQ: feed_line.New(),
	}

	// Start Workflow Channel IO in another goroutine
	// and send back the input as output (short circuit)
	go func() {
		for {
			select {
			case flu := <-w.InQ:
				w.OutQ <- flu
			}
		}
	}()

	return w
}
