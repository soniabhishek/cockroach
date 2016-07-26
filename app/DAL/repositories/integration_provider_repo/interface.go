package integration_provider_repo

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

//The public interface for IntegrationProviderRepo
type IIntegrationProviderRepo interface {
	GetById(Id uuid.UUID) (models.IntegrationProvider, error)
	GetByName(name string) (models.IntegrationProvider, error)
	Save(IntegrationProvider models.IntegrationProvider)
}

//--------------------------------------------------------MOCK----------------------------------------------------------------//

//Mock IntegrationProviderRepo. Used for testing other repos.
type Mock struct{}

//To ensure that it implements the interface
var _ IIntegrationProviderRepo = &Mock{}

func (i *Mock) GetById(Id uuid.UUID) (models.IntegrationProvider, error) {
	return getMockIP(), nil
}

func (i *Mock) GetByName(name string) (models.IntegrationProvider, error) {
	return getMockIP(), nil
}

func (i *Mock) Save(ip models.IntegrationProvider) {

}

func getMockIP() models.IntegrationProvider {

	return models.IntegrationProvider{
		ID:        uuid.NewV4(),
		Name:      "PayU",
		Label:     "payu",
		Website:   "payu.com",
		CreatedAt: pq.NullTime{time.Now(), true},
		UpdatedAt: pq.NullTime{time.Now(), true},
		LogoUrl:   sql.NullString{"https://easydigitaldownloads.com/wp-content/uploads/2013/12/payu-india-payment-gateway.png", true},
	}
}
