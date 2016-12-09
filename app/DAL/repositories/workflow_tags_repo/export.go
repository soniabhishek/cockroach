package workflow_tags_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkflowTagsRepo interface {
	Add([]models.WorkFlowTagAssociator) error
	Update([]models.WorkFlowTagAssociator) error
	Delete([]models.WorkFlowTagAssociator) error
	GetByWorkFlowId(id uuid.UUID) ([]models.WorkFlowTagAssociator, error)
	GetByProjectId(id uuid.UUID) ([]models.WorkFlowTagAssociator, error)
}

func New() IWorkflowTagsRepo {
	return &workflow_tags_repo{
		db: postgres.GetPostgresClient(),
	}
}
