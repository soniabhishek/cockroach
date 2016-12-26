package step_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepRepo_GetStartStep(t *testing.T) {

	pgClient := postgres.GetPostgresClient()

	var project models.Project
	err := pgClient.SelectOne(&project, `SELECT p.* from projects p LEFT OUTER JOIN work_flow w on w.project_id = p.id where w.id is null limit 1`)
	assert.NoError(t, err)

	workflow := models.WorkFlow{
		ID:        uuid.NewV4(),
		ProjectId: project.ID,
	}

	err = pgClient.Insert(&workflow)
	assert.NoError(t, err)

	startStep := models.Step{
		ID:         uuid.NewV4(),
		Type:       step_type.StartStep,
		WorkFlowId: workflow.ID,
	}
	err = pgClient.Insert(&startStep)
	assert.NoError(t, err)

	stepRepo := stepRepo{pgClient}
	startStepFromDb, err := stepRepo.GetStartStep(project.ID, "Test")
	assert.NoError(t, err)
	assert.Equal(t, startStep.ID, startStepFromDb.ID)

	_, err = pgClient.Delete(&startStep)
	assert.NoError(t, err)
	_, err = pgClient.Delete(&workflow)
	assert.NoError(t, err)

}
