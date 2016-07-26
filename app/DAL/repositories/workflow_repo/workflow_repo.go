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

func (wr *workflow_repo) Add(wf models.WorkFlow) error {
	return wr.db.Insert(&wf)
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
