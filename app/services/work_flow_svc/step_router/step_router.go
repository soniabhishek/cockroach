package step_router

import (
	"errors"

	"gitlab.com/playment-main/angel/app/DAL/repositories/step_router_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
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

	//Used for step router concurrency
	buffer chan uint
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

				// There is a question that adding to the
				// buffer should be inside or outside
				// the below go routine.
				//
				// Current implementation causes InQ to block if channel is full.
				// If put inside the below go routine InQ will become non blocking
				// but there is a chance of large number of go routines at a time
				sr.buffer <- 1

				go func() {

					defer func() {

						if r := recover(); r != nil {
							plog.Error("Router", errors.New("Panic occured in router"), r)
							*sr.routeTable[step.Manual] <- flu
						}

						// Free the buffer
						<-sr.buffer
					}()

					// Add workers here
					// Right now its like just one sync worker
					// i.e. if the below method takes 1 second & buffer size is 10
					// then the max speed of router processing will be 1 flu * buffer/second = 10 flu/second
					r, err := sr.getRoute(flu)
					if err != nil {
						plog.Error("Got route error , Sending it to manual", err, flu)
						*sr.routeTable[step.Manual] <- flu
					} else {
						counter.Print(flu, "router")
						*r <- flu
					}
				}()

			}
		}
	}()
}

func newStepRouter(concurrency int) stepRouter {
	return stepRouter{
		// Bigger feedLine since all the step servers
		// pushes flu to this one only
		InQ:           feed_line.NewBig(),
		ProcessedFluQ: feed_line.New(),
		buffer:        make(chan uint, concurrency),
		//stepRouterRepo: step_router_repo.Mock(),
	}
}

//--------------------------------------------------------------------------------//

type testInterface interface {
	GetNextStep(feed_line.FLU) (step.StepIdentifier, error)
}

var tt testInterface = &testStruct{
	stepRouterRepo: step_router_repo.Mock(),
}

type testStruct struct {
	// For Getting Next step
	stepRouterRepo step_router_repo.IStepRoutesRepo
}

func (t *testStruct) GetNextStep(flu feed_line.FLU) (step.StepIdentifier, error) {

	if l := len(flu.Trip); l > 0 {

		if currentStep := flu.Trip[l-1]; !currentStep.Success() {
			return step.Manual, nil
		}

	}

	routes, err := t.stepRouterRepo.GetRoutesByStepId(flu.StepId)
	if err != nil {
		return step.Manual, err
	}

	for _, route := range routes {
		correct, err := Logic(flu, models.LogicGate{ID: route.LogicGateId})
		if err != nil {
			return step.Manual, err
		} else if correct {
			return GetStepIdentifierFromStep(route.NextStepId), nil
		}
	}

	return step.Nil, nil
}

func GetStepIdentifierFromStep(stepId uuid.UUID) step.StepIdentifier {
	return step.CrowdSourcing
}
