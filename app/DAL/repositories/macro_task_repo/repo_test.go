package macro_task_repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

//Divide this test in setup & tear down
func TestGetMacroTask(t *testing.T) {

	pgClient := postgres.GetPostgresClient()
	mgo := clients.GetMongoClient()
	macroTaskRepo := macroTaskRepo{
		pg:  pgClient,
		mgo: clients.GetMongoClient(),
	}

	var pr models.Project
	err := pgClient.SelectOne(&pr, "select * from projects limit 1")
	ok := assert.NoError(t, err)
	if !ok {
		return
	}

	id := uuid.NewV4()
	m := models.MacroTask{
		ID:        id,
		Label:     "quality-check",
		Name:      "Quality Check",
		ProjectId: pr.ID,
		CreatorId: pr.CreatorId,
	}

	err = pgClient.Insert(&m)
	ok = assert.NoError(t, err)
	if !ok {
		return
	}
	defer func() {
		_, err = pgClient.Delete(&m)
		err1 := mgo.C("macro_tasks").RemoveId(m.ID)
		if err != nil {
			panic(err)
		}
		if err1 != nil {
			panic(err1)
		}
	}()

	err = macroTaskRepo.saveMgo(m)
	ok = assert.NoError(t, err)
	if !ok {
		return
	}

	mNew, err := macroTaskRepo.Get(m.ID)
	assert.NoError(t, err)
	ok = assert.EqualValues(t, m, mNew)
	if !ok {
		return
	}

	mNew, err = macroTaskRepo.getFromPG(m.ID)
	assert.NoError(t, err)
	ok = assert.Equal(t, m, mNew)
	if !ok {
		return
	}

	mNew, err = macroTaskRepo.getFromMgo(m.ID)
	assert.NoError(t, err)
	ok = assert.Equal(t, m, mNew)
	if !ok {
		return
	}
}

//Divide this test in setup & tear down
func TestAutoSaveMgo(t *testing.T) {
	pgClient := postgres.GetPostgresClient()
	mgo := clients.GetMongoClient()
	macroTaskRepo := macroTaskRepo{
		pg:  pgClient,
		mgo: clients.GetMongoClient(),
	}

	var pr models.Project
	err := pgClient.SelectOne(&pr, "select * from projects limit 1")
	ok := assert.NoError(t, err)
	if !ok {
		return
	}

	id := uuid.NewV4()
	m := models.MacroTask{
		ID:        id,
		Label:     "quality-check",
		Name:      "Quality Check",
		ProjectId: pr.ID,
		CreatorId: pr.CreatorId,
	}

	err = pgClient.Insert(&m)
	ok = assert.NoError(t, err)
	if !ok {
		return
	}
	defer func() {
		_, err = pgClient.Delete(&m)
		err1 := mgo.C("macro_tasks").RemoveId(m.ID)
		if err != nil {
			panic(err)
		}
		if err1 != nil {
			panic(err1)
		}
	}()

	_, err = macroTaskRepo.getFromMgo(m.ID)
	ok = assert.Error(t, err)
	if !ok {
		return
	}

	macroTaskRepo.Get(m.ID)

	//Waiting for macroTask to save in mongo on separate goroutine
	time.Sleep(time.Second * 1)

	mNew, err := macroTaskRepo.getFromMgo(m.ID)
	assert.NoError(t, err)
	ok = assert.Equal(t, m, mNew)
	if !ok {
		return
	}
}
