package clients_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IClientsRepo interface {
	GetByProjectId(projectId uuid.UUID) (models.Client, error)
	Add(models.Client) error
	Update(models.Client) error
	Delete(id uuid.UUID) error
}

func New() IClientsRepo {
	return &clientsRepo{
		Db: postgres.GetPostgresClient(),
	}
}
