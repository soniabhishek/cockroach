package step_configuration_repo

import "github.com/crowdflux/angel/app/DAL/clients/postgres"

func NewTransformationStepConfigurationRepo() ITransformationStepConfigurationRepo {
	return &transformationStepConfigRepo{Db: postgres.GetPostgresClient()}
}
