package projects_repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

//Divide this test in setup & tear down
func TestGetProject(t *testing.T) {

	pgClient := postgres.GetPostgresClient()
	projectsRepo := projectsRepo{
		pg:  pgClient,
		mgo: clients.GetMongoClient(),
	}

	var client models.Client
	pgClient.SelectOne(&client, "select * from clients limit 1")

	pr := models.Project{
		ID:        uuid.NewV4(),
		Name:      "Project 1",
		CreatorId: client.UserId,
		ClientId:  client.ID,
	}
	err := pgClient.Insert(&pr)
	ok := assert.NoError(t, err)
	if !ok {
		return
	}
	defer func() {
		pgClient.Delete(&pr)
	}()

	prNew, err := projectsRepo.Get(pr.ID)
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
	ok = assert.Equal(t, pr.ClientId, prNew.ClientId)
	if !ok {
		return
	}
}
