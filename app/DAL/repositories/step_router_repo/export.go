package step_router_repo

import "gitlab.com/playment-main/angel/app/DAL/clients/postgres"

func New() IStepRoutesRepo {
	return &stepRouteRepo{
		Db: postgres.GetPostgresClient(),
	}
}
