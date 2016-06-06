package step_router

import (
	"errors"

	"gitlab.com/playment-main/angel/app/DAL/repositories/step_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/step_router_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

type IRouteGetter interface {
	GetNextStep(feed_line.FLU) (models.Step, error)
	GetStartStep(flu feed_line.FLU) (models.Step, error)
}

type routeGetter struct {
	// For Getting Next step
	stepRouterRepo step_router_repo.IStepRoutesRepo
	stepRepo       step_repo.IStepRepo
}

func (r *routeGetter) GetStartStep(flu feed_line.FLU) (models.Step, error) {
	return r.stepRepo.GetStartStep(flu.ProjectId)

}

func (r *routeGetter) GetNextStep(flu feed_line.FLU) (models.Step, error) {

	var step models.Step

	routes, err := r.stepRouterRepo.GetRoutesByStepId(flu.StepId)
	if err != nil {
		return step, err
	}

	for _, route := range routes {
		correct, err := Logic(flu, models.LogicGate{ID: route.LogicGateId})
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
