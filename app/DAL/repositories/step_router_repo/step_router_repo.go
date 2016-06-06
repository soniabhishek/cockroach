package step_router_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type stepRouteRepo struct {
	Db repositories.IDatabase
}

var _ IStepRoutesRepo = (*stepRouteRepo)(nil)

const routesTable = "routes"

func (s *stepRouteRepo) GetRoutesByStepId(stepId uuid.UUID) (routes []models.Route, err error) {

	_, err = s.Db.Select(&routes, "select * from "+routesTable+" where step_id = $1", stepId)
	return
}

func (s *stepRouteRepo) Save(r models.Route) error {
	panic("not implemented")
	return nil
}
