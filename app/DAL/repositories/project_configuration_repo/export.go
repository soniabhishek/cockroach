package project_configuration_repo

import "gitlab.com/playment-main/angel/app/DAL/clients/postgres"

func New() IProjectConfigurationRepo {
	return &projectConfigurationRepo{
		db: postgres.GetPostgresClient(),
	}
}
