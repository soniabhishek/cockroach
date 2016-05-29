package step_router

import (
	"errors"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/crowdsourcing_step"
	"time"
)

type stepRouter struct {
	InQ           feed_line.FL
	ProcessedFluQ feed_line.FL
	routeTable    routeTable
}

type routeTable map[step.StepIdentifier]*feed_line.FL

type IConnector interface {

	// This method can be confusing
	Connect(routerIn *feed_line.FL) (routerOut *feed_line.FL)
}

func (sr *stepRouter) connectAll() {

	var connector IConnector = crowdsourcing_step.New()

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

func (sr *stepRouter) getRoute(flu models.FeedLineUnit) (*feed_line.FL, error) {

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

//--------------------------------------------------------------------------------//

type testInterface interface {
	GetNextStep(models.FeedLineUnit) (step.StepIdentifier, error)
}

var tt testInterface = testStruct{}

type testStruct struct {
}

func (testStruct) GetNextStep(flu models.FeedLineUnit) (step.StepIdentifier, error) {

	time.Sleep(time.Duration(2) * time.Second)

	return step.Nil, nil
}

//--------------------------------------------------------------------------------//
