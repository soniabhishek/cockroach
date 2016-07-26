package integration_provider_repo

import (
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/DAL/repositories/queries"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type integrationProviderRepo struct {
	Db repositories.IDatabase
}

var _ IIntegrationProviderRepo = &integrationProviderRepo{}

const integrationProviderTable string = "integration_providers"

func (i *integrationProviderRepo) GetById(id uuid.UUID) (models.IntegrationProvider, error) {

	var ip models.IntegrationProvider

	err := i.Db.SelectOne(&ip, queries.SelectById(integrationProviderTable), id)

	if err != nil {
		return ip, err
	}
	return ip, nil
}

func (i *integrationProviderRepo) GetByName(name string) (models.IntegrationProvider, error) {

	var ip models.IntegrationProvider

	err := i.Db.SelectOne(&ip, queries.SelectByName(integrationProviderTable), name)
	if err != nil {
		return ip, err
	}
	return ip, nil
}

func (i *integrationProviderRepo) Save(ip models.IntegrationProvider) {

}
