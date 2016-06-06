package step_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IStepRepo interface {
	GetById(id uuid.UUID) (models.Step, error)
	GetStartStep(projectId uuid.UUID) (models.Step, error)
	GetEndStep(projectId uuid.UUID) (models.Step, error)
}

func New() IStepRepo {
	return &stepRepo{
		Db: postgres.GetPostgresClient(),
	}
}
