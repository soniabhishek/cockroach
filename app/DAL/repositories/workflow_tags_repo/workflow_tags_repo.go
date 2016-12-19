package workflow_tags_repo

import (
	"database/sql"
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"time"
)

type workflow_tags_repo struct {
	Db repositories.IDatabase
}

var _ IWorkflowTagsRepo = &workflow_tags_repo{}

func (wtr *workflow_tags_repo) Add(wfTags []models.WorkFlowTagAssociator) error {
	var tagsInterface []interface{} = make([]interface{}, len(wfTags))
	for i, _ := range wfTags {
		wfTags[i].CreatedAt = pq.NullTime{time.Now(), true}
		tagsInterface[i] = &wfTags[i]
	}

	return wtr.Db.Insert(tagsInterface...)
}

func (wtr *workflow_tags_repo) Update(wfTags []models.WorkFlowTagAssociator) error {
	var tagsInterface []interface{} = make([]interface{}, len(wfTags))
	for i, _ := range wfTags {
		wfTags[i].CreatedAt = pq.NullTime{time.Now(), true}
		tagsInterface[i] = &wfTags[i]
	}
	_, err := wtr.Db.Update(tagsInterface...)
	return err
}

func (wtr *workflow_tags_repo) Delete(wfTags []models.WorkFlowTagAssociator) error {
	var tagsInterface []interface{} = make([]interface{}, len(wfTags))
	for i, _ := range wfTags {
		tagsInterface[i] = &wfTags[i]
	}
	_, err := wtr.Db.Delete(tagsInterface...)
	return err
}
func (wtr *workflow_tags_repo) GetByWorkFlowId(id uuid.UUID) (wfTags []models.WorkFlowTagAssociator, err error) {
	_, err = wtr.Db.Select(&wfTags, `select * from work_flow_tag_associators where work_flow_id = $1  `, id)

	if err == sql.ErrNoRows || len(wfTags) == 0 {
		err = ErrWorkflowTagsNotFound
	}

	return
}
func (wtr *workflow_tags_repo) GetByProjectId(id uuid.UUID) (wfTags []models.WorkFlowTagAssociator, err error) {
	_, err = wtr.Db.Select(&wfTags, `select * from work_flow_tag_associators where project_id = $1  `, id)
	return
}
