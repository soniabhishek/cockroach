package step_repo

import (
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type stepRepo struct {
	Db repositories.IDatabase
}

const stepTable = "step"

var _ IStepRepo = &stepRepo{}

func (s *stepRepo) GetById(id uuid.UUID) (step models.Step, err error) {
	err = s.Db.SelectById(&step, id)
	return
}

func (s *stepRepo) GetStepsByWorkflowId(workFlowId uuid.UUID) (steps []models.Step, err error) {

	_, err = s.Db.Select(&steps, `select * from step where work_flow_id = $1 `, workFlowId.String())
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

func (s *stepRepo) AddMany(steps []models.Step) (err error) {
	var stepsInterface []interface{} = make([]interface{}, len(steps))
	for i, _ := range steps {
		stepsInterface[i] = &steps[i]
	}

	err = s.Db.Insert(stepsInterface...)
	return
}

func (s *stepRepo) UpdateMany(steps []models.Step) (result int64, err error) {
	var stepsInterface []interface{} = make([]interface{}, len(steps))
	for i, _ := range steps {
		stepsInterface[i] = &steps[i]
	}

	result, err = s.Db.Update(stepsInterface...)
	return
}
func (s *stepRepo) DeleteMany(steps []models.Step) (result int64, err error) {
	var stepsInterface []interface{} = make([]interface{}, len(steps))
	for i, _ := range steps {
		stepsInterface[i] = &steps[i]
	}

	result, err = s.Db.Delete(stepsInterface...)
	return
}
