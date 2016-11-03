package router_svc

import (
	"errors"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_router_repo"
	"github.com/crowdflux/angel/app/models"
)

type IRouteGetter interface {
	GetNextStep(feed_line.FLU) (models.Step, error)
}

type routeGetter struct {
	// For Getting Next step
	stepRouterRepo step_router_repo.IStepRoutesRepo
	stepRepo       step_repo.IStepRepo
}

func (r *routeGetter) GetNextStep(flu feed_line.FLU) (models.Step, error) {

	var step models.Step

	routes, err := r.stepRouterRepo.GetRoutesByStepId(flu.StepId)
	if err != nil {
		return step, err
	}

	for _, route := range routes {

		var logicGate models.LogicGate
		logicGate.InputTemplate = route.Config
		correct, err := Logic(flu, logicGate)
		if err != nil {
			return step, err
		} else if correct {
			stp, err := r.stepRepo.GetById(route.NextStepId)
			return stp, err
		}
	}

	return step, errors.New("no matching route")
}

func newRouteGetter() IRouteGetter {
	return &routeGetter{
		stepRouterRepo: step_router_repo.New(),
		stepRepo:       step_repo.New(),
	}
}
