package integration_provider_repo

import (
	"gitlab.com/playment-main/support/app/services/data_access_svc/clients"
	"gitlab.com/playment-main/support/app/services/data_access_svc/repositories"
)

func New() IIntegrationProviderRepo {
	return NewCustom(clients.GetPostgresClient())
}

func NewCustom(dbInterface repositories.IDatabase) IIntegrationProviderRepo {
	return &integrationProviderRepo{
		Db: dbInterface,
	}
}
