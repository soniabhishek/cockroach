package workflow_tags_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkflowTagsRepo interface {
	Add([]models.WorkFlowTagAssociators) error
	Update([]models.WorkFlowTagAssociators) error
	Delete([]models.WorkFlowTagAssociators) error
	GetByWorkFlowId(id uuid.UUID) ([]models.WorkFlowTagAssociators, error)
}

func New() IWorkflowTagsRepo {
	return &workflow_tags_repo{
		db: postgres.GetPostgresClient(),
	}
}
