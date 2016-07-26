package integration_provider_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/DAL/repositories"
)

func New() IIntegrationProviderRepo {
	return NewCustom(postgres.GetPostgresClient())
}

func NewCustom(dbInterface repositories.IDatabase) IIntegrationProviderRepo {
	return &integrationProviderRepo{
		Db: dbInterface,
	}
}
