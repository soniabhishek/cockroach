package work_flow

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/services/work_flow_svc/counter"
	"github.com/crowdflux/angel/app/services/work_flow_svc/router_svc"
)

type WorkFlow struct {
	InQ  feed_line.Fl
	OutQ feed_line.Fl
}

// Creates a new workFlow instance
// Making it private because making it public can be dangerous
func newStdWorkFlow() WorkFlow {

	//create new instance
	w := WorkFlow{
		InQ:  feed_line.New("workflow-in"),
		OutQ: feed_line.New("workflow-out"),
	}

	// Start Workflow Channel IO in another goroutine
	go func() {

		inputQueue := w.InQ.Receiver()
		outputQueue := router_svc.StdStepRouter.ProcessedFluQ.Receiver()

		for {
			select {

			case flu := <-inputQueue:
				counter.Print(flu, "workflow in")
				router_svc.StdStepRouter.InQ.Push(flu)

			case flu := <-outputQueue:
				counter.Print(flu, "workflow out")
				w.OutQ.Push(flu)
			}
		}
	}()

	return w
}

// Exposing a StdWorkFlow instance
var StdWorkFlow = newStdWorkFlow()

//var StdWorkFlow = newShortCircuit()

func NewShortCircuit() WorkFlow {
	//create new instance
	w := WorkFlow{
		InQ:  feed_line.New("workflow-in2123"),
		OutQ: feed_line.New("workflow-out2123"),
	}

	// Start Workflow Channel IO in another goroutine
	// and send back the input as output (short circuit)
	go func() {
		for flu := range w.InQ.Receiver() {
			w.OutQ.Push(flu)
		}
	}()

	return w
}
