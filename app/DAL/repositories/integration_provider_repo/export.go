package integration_provider_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/DAL/repositories"
)

func New() IIntegrationProviderRepo {
	return NewCustom(postgres.GetPostgresClient())
}

func NewCustom(dbInterface repositories.IDatabase) IIntegrationProviderRepo {
	return &integrationProviderRepo{
		Db: dbInterface,
	}
}
