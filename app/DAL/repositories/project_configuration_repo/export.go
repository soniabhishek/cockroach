package project_configuration_repo

import "github.com/crowdflux/angel/app/DAL/clients/postgres"

func New() IProjectConfigurationRepo {
	return &projectConfigurationRepo{
		db: postgres.GetPostgresClient(),
	}
}
