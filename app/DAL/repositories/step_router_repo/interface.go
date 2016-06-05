package step_router_repo

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"sync"
)

type IStepRoutesRepo interface {
	GetRoutesByStepId(stepId uuid.UUID) ([]models.Route, error)
	Save(models.Route) error
}

//========================================================================================================================================================//

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

func (i *inMemStepRouteRepo) Save(r models.Route) error {
	i.Lock()
	defer i.Unlock()
	i.stepRoutes[r.ID] = r
	return nil
}

func Mock() IStepRoutesRepo {
	return &inMemStepRouteRepo{
		stepRoutes: make(map[uuid.UUID]models.Route),
	}
}
