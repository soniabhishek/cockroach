package projects_repo

import (
	"testing"

	"time"

	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

//Divide this test in setup & tear down
func TestGetProject(t *testing.T) {

	pgClient := postgres.GetPostgresClient()
	projectsRepo := projectsRepo{
		pg:  pgClient,
		mgo: clients.GetMongoClient(),
	}

	var client models.Client
	err := pgClient.SelectOne(&client, "select * from clients limit 1")
	assert.NoError(t, err)

	pr := models.Project{
		ID:        uuid.NewV4(),
		Name:      "Project 1",
		CreatorId: client.UserId,
		ClientId:  client.ID,
	}
	err = pgClient.Insert(&pr)
	ok := assert.NoError(t, err)
	if !ok {
		return
	}
	defer func() {
		pgClient.Delete(&pr)
	}()

	prNew, err := projectsRepo.GetById(pr.ID)
	// Give the above line some time
	// to save data in mongo since it will
	// run in another go routine
	time.Sleep(time.Duration(100) * time.Millisecond)

	assert.NoError(t, err)
	ok = assert.EqualValues(t, pr, prNew)
	if !ok {
		return
	}

	prNew, err = projectsRepo.getFromPG(pr.ID)
	assert.NoError(t, err)
	ok = assert.Equal(t, pr, prNew)
	if !ok {
		return
	}

	prNew, err = projectsRepo.getFromMgo(pr.ID)
	assert.NoError(t, err)
	assert.Equal(t, pr.ClientId, prNew.ClientId)
}
