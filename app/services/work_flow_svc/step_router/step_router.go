package step_router

import (
	"errors"
	"time"

	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/crowdsourcing_step"
)

type routeTable map[step.StepIdentifier]*feed_line.Fl

type IConnector interface {

	// This method can be confusing
	Connect(routerIn *feed_line.Fl) (routerOut *feed_line.Fl)
}

type stepRouter struct {
	InQ           feed_line.Fl
	ProcessedFluQ feed_line.Fl
	routeTable    routeTable
}

func (sr *stepRouter) connectAll() {

	var connector IConnector = crowdsourcing_step.StdCrowdSourcingStep

	sr.routeTable = routeTable{

		step.CrowdSourcing:    connector.Connect(&sr.InQ),
		step.InternalSourcing: connector.Connect(&sr.InQ),
		step.Manual:           connector.Connect(&sr.InQ),
		step.Transformation:   connector.Connect(&sr.InQ),
		step.Algorithm:        connector.Connect(&sr.InQ),
		step.Bifurcation:      connector.Connect(&sr.InQ),
		step.Unification:      connector.Connect(&sr.InQ),

		// Special case
		// Map route end to Processed Flu Queue (ProcessedFluQ)
		step.Nil: &sr.ProcessedFluQ,
	}
}

func (sr *stepRouter) getRoute(flu feed_line.FLU) (*feed_line.Fl, error) {

	next, err := tt.GetNextStep(flu)
	if err != nil {
		return nil, err
	}

	if route, ok := sr.routeTable[next]; ok {
		return route, nil
	} else {
		return nil, errors.New("Step not found")
	}
}

func (sr *stepRouter) start() {

	// Start listening for incoming flus from InQ channel
	// in another goroutine & route it to its exact step
	go func() {
		for {
			select {
			case flu := <-sr.InQ:

				// Add workers here
				// Right now its like just one sync worker
				// i.e. if the below method takes 1 second
				// then the max speed of router processing will be 1 flu/second
				r, _ := sr.getRoute(flu)
				flu.Step = "router"
				counter.Print(flu)
				*r <- flu
			}
		}
	}()
}

func newStepRouter() stepRouter {
	return stepRouter{
		// Bigger feedLine since all the step servers
		// pushes flu to this one only
		InQ:           feed_line.NewBig(),
		ProcessedFluQ: feed_line.New(),
	}
}

//--------------------------------------------------------------------------------//

type testInterface interface {
	GetNextStep(feed_line.FLU) (step.StepIdentifier, error)
}

var tt testInterface = testStruct{}

type testStruct struct {
}

func (testStruct) GetNextStep(flu feed_line.FLU) (step.StepIdentifier, error) {

	time.Sleep(time.Duration(200) * time.Millisecond)

	return step.Nil, nil
}

//--------------------------------------------------------------------------------//
