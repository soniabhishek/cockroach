package feed_line_repo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
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
