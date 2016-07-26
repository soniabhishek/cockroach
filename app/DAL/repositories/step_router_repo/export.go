package step_router_repo

import "github.com/crowdflux/angel/app/DAL/clients/postgres"

func New() IStepRoutesRepo {
	return &stepRouteRepo{
		Db: postgres.GetPostgresClient(),
	}
}
