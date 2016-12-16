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
	"strconv"
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
	flus = append(flus, models.FeedLineUnit{ID: v1, Build: models.JsonF{"1": "One"}, ProjectId: projectId, StepId: step.ID, MasterId: v1})
	flus = append(flus, models.FeedLineUnit{ID: v2, Build: models.JsonF{"2": "Two"}, ProjectId: projectId, StepId: step.ID, MasterId: v2})
	flus = append(flus, models.FeedLineUnit{ID: v3, Build: models.JsonF{"3": "Three"}, ProjectId: projectId, StepId: step.ID, MasterId: v3})
	flus = append(flus, models.FeedLineUnit{ID: v4, Build: models.JsonF{"4": "Four"}, ProjectId: projectId, StepId: step.ID, MasterId: v4})

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

func TestBulkFluInsert(t *testing.T) {
	t1, _ := uuid.FromString("214a98dd-e95a-4d23-99fb-c76e836afe62")
	t2, _ := uuid.FromString("c9177ac2-c47b-4968-9b11-96c3d6c84712")
	t3, _ := uuid.FromString("e436b759-fa12-4113-9f39-0f5ab934f2b6")
	t4, _ := uuid.FromString("3e4e2c42-82ec-4877-b123-115f8e5858a4")
	t5, _ := uuid.FromString("64610638-02b3-41a2-b891-07caeae68493")
	idmap := []uuid.UUID{t1, t2, t3, t4, t5}
	flus := make([]models.FeedLineUnit, 0)
	for i := 0; i < 5; i++ {
		flus = append(flus, models.FeedLineUnit{ID: idmap[i], Build: models.JsonF{strconv.Itoa(i): "One"}})
	}
	conn := NewInputQueue()
	err := conn.BulkAdd(flus)
	assert.NoError(t, err)

}
