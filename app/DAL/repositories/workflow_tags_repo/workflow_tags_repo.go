package workflow_tags_repo

import (
	"database/sql"
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type workflow_tags_repo struct {
	db repositories.IDatabase
}

var _ IWorkflowTagsRepo = &workflow_tags_repo{}

func (wtr *workflow_tags_repo) Add(wfTags []models.WorkFlowTagAssociators) error {
	var tagsInterface []interface{} = make([]interface{}, len(wfTags))
	for i, _ := range wfTags {
		tagsInterface[i] = &wfTags[i]
	}

	return wtr.db.Insert(tagsInterface...)
}

func (wtr *workflow_tags_repo) Update(wfTags []models.WorkFlowTagAssociators) error {
	var tagsInterface []interface{} = make([]interface{}, len(wfTags))
	for i, _ := range wfTags {
		tagsInterface[i] = &wfTags[i]
	}
	_, err := wtr.db.Update(tagsInterface...)
	return err
}

func (wtr *workflow_tags_repo) Delete(wfTags []models.WorkFlowTagAssociators) error {
	var tagsInterface []interface{} = make([]interface{}, len(wfTags))
	for i, _ := range wfTags {
		tagsInterface[i] = &wfTags[i]
	}
	_, err := wtr.db.Delete(tagsInterface...)
	return err
}
func (wtr *workflow_tags_repo) GetByWorkFlowId(id uuid.UUID) (wfTags []models.WorkFlowTagAssociators, err error) {
	_, err = wtr.db.Select(&wfTags, `select * from work_flow_tag_associators where work_flow_id = $1  `, id)

	if err == sql.ErrNoRows || len(wfTags) == 0 {
		err = ErrWorkflowTagsNotFound
	}

	return
}
