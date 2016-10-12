package step_router_repo

import (
	"sync"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type inMemStepRouteRepo struct {
	sync.RWMutex
	stepRoutes map[uuid.UUID]models.Route
}

var _ IStepRoutesRepo = (*inMemStepRouteRepo)(nil)

func (i *inMemStepRouteRepo) GetRoutesByStepId(stepId uuid.UUID) (routes []models.Route, err error) {
	i.RLock()
	defer i.RUnlock()

	for _, route := range i.stepRoutes {
		if route.StepId == stepId {
			routes = append(routes, route)
		}
	}

	if len(routes) == 0 {
		return routes, ErrRouteNotFound
	}
	return routes, nil
}
func (i *inMemStepRouteRepo) GetRoutesWithLogicByStepId(stepId uuid.UUID) (routes []models.RouteWithLogicGate, err error) {
	panic("not implemented")
	return routes, nil
}
func (i *inMemStepRouteRepo) GetRoutesByWorkFlowId(stepId uuid.UUID) (routes []models.Route, err error) {
	return
}

func (i *inMemStepRouteRepo) Save(r models.Route) error {
	i.Lock()
	defer i.Unlock()
	i.stepRoutes[r.ID] = r
	return nil
}
func (i *inMemStepRouteRepo) AddMany(routes []models.Route) error {
	return nil
}
func (i *inMemStepRouteRepo) UpdateMany(routes []models.Route) (result int64, err error) {
	return
}
func (i *inMemStepRouteRepo) DeleteMany(routes []models.Route) (result int64, err error) {
	return
}
func Mock() IStepRoutesRepo {
	return &inMemStepRouteRepo{
		stepRoutes: make(map[uuid.UUID]models.Route),
	}
}
