package step_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IStepRepo interface {
	GetById(id uuid.UUID) (models.Step, error)
	GetStartStep(projectId uuid.UUID) (models.Step, error)
	GetEndStep(projectId uuid.UUID) (models.Step, error)
	GetStepsByWorkflowId(id uuid.UUID) ([]models.Step, error)
	AddMany([]models.Step) error
	UpdateMany([]models.Step) (int64, error)
	DeleteMany([]models.Step) (int64, error)
}

func New() IStepRepo {
	return &stepRepo{
		Db: postgres.GetPostgresClient(),
	}
}
