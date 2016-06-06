package step_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type stepRepo struct {
	Db repositories.IDatabase
}

const stepTable = "step"

func (s *stepRepo) GetById(id uuid.UUID) (step models.Step, err error) {
	err = s.Db.SelectById(&step, id)
	return
}

func (s *stepRepo) GetStartStep(projectId uuid.UUID) (step models.Step, err error) {

	err = s.Db.SelectOne(&step, `
	select s.* from step s
	inner join work_flow w on w.id = s.work_flow_id
	where w.project_id = $1 and is_start is true`, projectId.String())
	return
}

func (s *stepRepo) GetEndStep(projectId uuid.UUID) (step models.Step, err error) {
	err = s.Db.SelectOne(&step, `
	select s.* from step s
	inner join work_flow w on w.id = s.work_flow_id
	where w.project_id = $1 and is_end is true`, projectId.String())
	return
}
