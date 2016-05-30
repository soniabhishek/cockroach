package flu_project_service

import "gitlab.com/playment-main/angel/app/DAL/clients/postgres"

func New() IFluProjectServiceRepo {
	return &fluProjectService{
		db: postgres.GetPostgresClient(),
	}
}
