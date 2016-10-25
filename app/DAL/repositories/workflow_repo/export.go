package workflow_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkflowRepo interface {
	Add(*models.WorkFlow) error
	Update(models.WorkFlow) error
	Delete(uuid.UUID) error
	GetById(id uuid.UUID) (models.WorkFlow, error)
	GetWorkFlowByProjectIdAndTag(projectId uuid.UUID, tag string)  (models.WorkFlow, error)
	GetWorkFlowsByProjectId(projectId uuid.UUID) ([]models.WorkFlow, error)
	IfIdExist(uuid.UUID) (bool, error)
}

func New() IWorkflowRepo {
	return &workflow_repo{
		db: postgres.GetPostgresClient(),
	}
}
