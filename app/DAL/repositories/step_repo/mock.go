package step_repo

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type stepRepoMock struct {
}

var _ IStepRepo = &stepRepoMock{}

func (s *stepRepoMock) GetById(id uuid.UUID) (step models.Step, err error) { return }
func (s *stepRepoMock) GetStartStep(projectId uuid.UUID, tag string) (step models.Step, err error) {
	return
}
func (s *stepRepoMock) GetStartStepOrDefault(projectId uuid.UUID, tag string) (step models.Step, err error) {
	return
}
func (s *stepRepoMock) GetEndStep(projectId uuid.UUID) (step models.Step, err error)       { return }
func (s *stepRepoMock) GetStepsByWorkflowId(id uuid.UUID) (steps []models.Step, err error) { return }
func (s *stepRepoMock) AddMany([]models.Step) (err error)                                  { return }
func (s *stepRepoMock) UpdateMany([]models.Step) (i int64, err error)                      { return }
func (s *stepRepoMock) DeleteMany([]models.Step) (i int64, err error)                      { return }

func Mock() IStepRepo {
	return &stepRepoMock{}
}
