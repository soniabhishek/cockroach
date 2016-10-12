package router_svc

import (
	"errors"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step/crowdsourcing_step_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step/manual_step_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step/transformation_step_csv"
)

type routeTable map[step_type.StepType]*feed_line.Fl

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

	// For Getting routes from storage
	routeGetter IRouteGetter

	// For saving flu
	fluRepo feed_line_repo.IFluRepo
}

func (sr *stepRouter) connectAll() {

	var crowdSourcingConn IConnector = crowdsourcing_step_svc.StdCrowdSourcingStep
	var manualStepConn IConnector = manual_step_svc.StdManualStep
	var transformationStepConn IConnector = transformation_step_svc.StdTransformationStep

	sr.routeTable = routeTable{

		step_type.CrowdSourcing:    crowdSourcingConn.Connect(&sr.InQ),
		step_type.InternalSourcing: crowdSourcingConn.Connect(&sr.InQ),
		step_type.Manual:           manualStepConn.Connect(&sr.InQ),
		step_type.Transformation:   transformationStepConn.Connect(&sr.InQ),
		step_type.Algorithm:        manualStepConn.Connect(&sr.InQ),
		step_type.Bifurcation:      manualStepConn.Connect(&sr.InQ),
		step_type.Unification:      manualStepConn.Connect(&sr.InQ),
		step_type.Error:            manualStepConn.Connect(&sr.InQ),

		// Special case
		// Map route end to Processed Flu Queue (ProcessedFluQ)
		step_type.Gateway: &sr.ProcessedFluQ,
	}
}

func (sr *stepRouter) getRoute(flu *feed_line.FLU) (route *feed_line.Fl) {

	var nextStep models.Step
	var err error

	// If flu's step id is nil
	// then its a new flu directly from outside
	// Get that step or send it to error step
	if flu.StepId == uuid.Nil {
		nextStep, err = sr.routeGetter.GetStartStep(*flu)
		if err != nil {
			plog.Error("Router", err, "error occured while getting start step")
			return sr.routeTable[step_type.Error]
		}

	} else {

		nextStep, err = sr.routeGetter.GetNextStep(*flu)
		if err != nil {
			// Error getting the next step
			plog.Error("Router", err, "error while getting evaluating logics in get route")
			return sr.routeTable[step_type.Error]
		}
	}

	if route, ok := sr.routeTable[nextStep.Type]; ok {

		// save the flu state change
		flu.StepId = nextStep.ID
		err := sr.fluRepo.Update(flu.FeedLineUnit)
		if err != nil {
			plog.Error("Router", err, "error occured while saving flu in router")
			return sr.routeTable[step_type.Error]
		}

		// Return the correct route
		plog.Info("router", "sending flu to ", nextStep.ID.String(), " Type", nextStep.Type)
		return route
	} else {

		// Route not found in the route table
		plog.Error("Router", errors.New("route not found in route table"))
		return sr.routeTable[step_type.Error]
	}
}

func newStepRouter(concurrency int) stepRouter {
	return stepRouter{
		// Bigger feedLine since all the step servers
		// pushes flu to this one only
		InQ:           feed_line.New("router-in"),
		ProcessedFluQ: feed_line.New("router-out-processed-flu"),
		buffer:        make(chan uint, concurrency),
		routeGetter:   newRouteGetter(),
		fluRepo:       feed_line_repo.New(),
	}
}
