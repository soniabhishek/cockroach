package feed_line_repo

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/stretchr/testify/assert"
)

func TestFeedLineGet(t *testing.T) {
	pgCient := postgres.GetPostgresClient()

	var fluId uuid.UUID

	err := pgCient.SelectOne(&fluId, "select id from feed_line limit 1")
	assert.NoError(t, err)
	fmt.Println(fluId)

	fluRepo := &fluRepo{
		Db: pgCient,
	}
	flu, err := fluRepo.GetById(fluId)
	assert.NoError(t, err)
	assert.Equal(t, fluId, flu.ID)
	fmt.Println(flu)

}

func TestBulkFluBuildUpdate(t *testing.T) {

	pgCient := postgres.GetPostgresClient()
	var projectId uuid.UUID
	var step models.Step

	pgCient.SelectOne(&projectId, "select project_id from feed_line limit 1")
	pgCient.SelectOne(&step, "select * from step limit 1")

	flus := make([]models.FeedLineUnit, 0)

	v1 := uuid.NewV4()
	v2 := uuid.NewV4()
	v3 := uuid.NewV4()
	v4 := uuid.NewV4()
	flus = append(flus, models.FeedLineUnit{ID: v1, Build: models.JsonF{"1": "One"}, ProjectId: projectId, StepId: step.ID})
	flus = append(flus, models.FeedLineUnit{ID: v2, Build: models.JsonF{"2": "Two"}, ProjectId: projectId, StepId: step.ID})
	flus = append(flus, models.FeedLineUnit{ID: v3, Build: models.JsonF{"3": "Three"}, ProjectId: projectId, StepId: step.ID})
	flus = append(flus, models.FeedLineUnit{ID: v4, Build: models.JsonF{"4": "Four"}, ProjectId: projectId, StepId: step.ID})

	flRepo := fluRepo{
		Db:       postgres.GetPostgresClient(),
		stepRepo: step_repo.New(),
	}

	err := flRepo.Add(flus[0])
	assert.NoError(t, err)

	_, err = flRepo.BulkFluBuildUpdateByStepType(flus, step.Type)
	assert.Error(t, err)

	updatedFlus, err := flRepo.BulkFluBuildUpdateByStepType(flus[0:1], step.Type)
	plog.Info("asd", updatedFlus)
	assert.NoError(t, err)

	flu, err := flRepo.GetById(v1)
	assert.NoError(t, err)
	assert.True(t, flu.UpdatedAt.Valid)

	pgCient.Delete(&flu)

}
