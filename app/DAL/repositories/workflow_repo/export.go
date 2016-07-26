package workflow_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
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
