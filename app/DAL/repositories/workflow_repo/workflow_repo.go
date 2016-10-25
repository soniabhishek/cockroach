package workflow_repo

import (
	"errors"
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type workflow_repo struct {
	db repositories.IDatabase
}

var _ IWorkflowRepo = &workflow_repo{}

func (wr *workflow_repo) Add(wf *models.WorkFlow) error {
	return wr.db.Insert(wf)
}

func (wr *workflow_repo) Update(wf models.WorkFlow) error {
	_, err := wr.db.Update(&wf)
	return err
}

func (wr *workflow_repo) Delete(id uuid.UUID) error {
	query := fmt.Sprintf(`delete from work_flow where id='%v'::uuid`, id)
	res, err := wr.db.Exec(query)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		err = errors.New("Could not delete WorkFlow with ID [" + id.String() + "]")
	}
	return err
}
func (wr *workflow_repo) GetById(id uuid.UUID) (wf models.WorkFlow, err error) {
	wf = models.WorkFlow{}
	err = wr.db.SelectById(&wf, id)
	return
}

func (wr *workflow_repo) GetWorkFlowByProjectIdAndTag(projectId uuid.UUID, tag string) (workFlow models.WorkFlow, err error) {
	err = wr.db.SelectOne(&workFlow, `select * from work_flow where project_id = $1 and tag = $2 `, projectId.String(), tag)
	return

}

func (wr *workflow_repo) GetWorkFlowsByProjectId(projectId uuid.UUID) (workFlows []models.WorkFlow, err error) {
	_, err = wr.db.Select(&workFlows, `select * from work_flow where project_id = $1`, projectId.String())
	fmt.Println(workFlows)
	return
}


func (i *workflow_repo) IfIdExist(id uuid.UUID) (ifExist bool, err error) {
	err = i.db.SelectOne(&ifExist, `select exists(select 1 from work_flow where id=$1)`, id)
	if err != nil {
		return
	}
	return
}
