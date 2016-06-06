package clients_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IClientsRepo interface {
	GetByProjectId(projectId uuid.UUID) (models.Client, error)
}

func New() IClientsRepo {
	return &clientsRepo{
		Db: postgres.GetPostgresClient(),
	}
}
