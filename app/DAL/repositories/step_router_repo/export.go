package step_router_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IStepRoutesRepo interface {
	GetRoutesByStepId(stepId uuid.UUID) ([]models.Route, error)
	GetRoutesWithLogicByStepId(stepId uuid.UUID) ([]models.RouteWithLogicGate, error)
	Save(models.Route) error
	UpdateMany([]models.Route) (int64, error)
	DeleteMany([]models.Route) (int64, error)
	AddMany([]models.Route) error
	GetRoutesByWorkFlowId(workFlowId uuid.UUID) ([]models.Route, error)
}

func New() IStepRoutesRepo {
	return &stepRouteRepo{
		Db: postgres.GetPostgresClient(),
	}
}
