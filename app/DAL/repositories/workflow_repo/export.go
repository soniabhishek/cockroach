package workflow_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IWorkflowRepo interface {
	Add(models.WorkFlow) error
	Update(models.WorkFlow) error
	Delete(uuid.UUID) error
}

func New() IWorkflowRepo {
	return &workflow_repo{
		db: postgres.GetPostgresClient(),
	}
}
