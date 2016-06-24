package feed_line_repo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"testing"
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
	flus := make([]models.FeedLineUnit, 0)

	v1, _ := uuid.FromString("9d2ecb8b-84ce-4d92-a01c-fa3ba71e1b61")
	v2, _ := uuid.FromString("507b699a-192f-40e1-9d40-bb7689d41962")
	v3, _ := uuid.FromString("5525b18d-3b5b-4d57-9e88-bddbfd1653e3")
	v4, _ := uuid.FromString("6179be61-4325-429a-8e51-e8da0d8672a5")
	flus = append(flus, models.FeedLineUnit{ID: v1, Build: models.JsonFake{"1": "One"}})
	flus = append(flus, models.FeedLineUnit{ID: v2, Build: models.JsonFake{"2": "Two"}})
	flus = append(flus, models.FeedLineUnit{ID: v3, Build: models.JsonFake{"3": "Three"}})
	flus = append(flus, models.FeedLineUnit{ID: v4, Build: models.JsonFake{"4": "Four"}})

	flRepo := New()
	err := flRepo.BulkFluBuildUpdate(flus)
	fmt.Println(err)
}
