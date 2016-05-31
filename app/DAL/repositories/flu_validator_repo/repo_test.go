package flu_validator_repo

import (
	"testing"
	"time"

	"os"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

var testMacroTask models.MacroTask

func TestInsertGetDelete(t *testing.T) {

	id := uuid.NewV4()

	validator := models.FLUValidator{
		ID:          id,
		FieldName:   "brand",
		Type:        "STRING",
		IsMandatory: true,
		MacroTaskId: testMacroTask.ID,
		CreatedAt:   pq.NullTime{time.Now(), true},
	}

	r := fluValidatorRepo{
		db: postgres.GetPostgresClient(),
	}

	err := r.insertMany([]models.FLUValidator{validator})

	e := assert.NoError(t, err, "Error occured while saving", err)

	if !e {
		return
	}

	_, err = r.get(id)

	e = assert.NoError(t, err, "Error occured while getting", err)

	if !e {
		return
	}

	err = r.deleteMany([]models.FLUValidator{validator})
	assert.NoError(t, err, "Error occured while deleting", err)
}

func TestSaveExisting(t *testing.T) {
	id := uuid.NewV4()

	validator := models.FLUValidator{
		ID:          id,
		FieldName:   "brand",
		Type:        "STRING",
		IsMandatory: true,
		MacroTaskId: testMacroTask.ID,
		CreatedAt:   pq.NullTime{time.Now(), true},
	}

	r := fluValidatorRepo{
		db: postgres.GetPostgresClient(),
	}

	err := r.insertMany([]models.FLUValidator{validator})

	e := assert.NoError(t, err, "Error occured while inserting", err)
	if !e {
		return
	}
	validator.IsMandatory = false
	err = r.Save(&validator)
	assert.NoError(t, err, "Error Occurred")

	flus, err := r.get(validator.ID)

	assert.NoError(t, err)
	assert.Equal(t, flus.IsMandatory, false)

	err = r.deleteMany([]models.FLUValidator{validator})
	assert.NoError(t, err, "Error occured while deleting", err)

}

func TestSaveNew(t *testing.T) {
	id := uuid.NewV4()

	validator := models.FLUValidator{
		ID:          id,
		FieldName:   "brand",
		Type:        "STRING",
		IsMandatory: true,
		MacroTaskId: testMacroTask.ID,
		CreatedAt:   pq.NullTime{time.Now(), true},
	}

	r := fluValidatorRepo{
		db: postgres.GetPostgresClient(),
	}

	err := r.Save(&validator)
	assert.NoError(t, err, "Error occured while saving")

	flus, err := r.get(validator.ID)

	assert.NoError(t, err)
	assert.Equal(t, flus.IsMandatory, true)

	err = r.deleteMany([]models.FLUValidator{validator})
	assert.NoError(t, err, "Error occured while deleting", err)

}

func TestMain(m *testing.M) {
	// your func
	setup()

	retCode := m.Run()

	// your func
	teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

func setup() {

	//Load any macro_task from db
	//Make sure you have a macro_task in db
	err := postgres.GetPostgresClient().SelectOne(&testMacroTask, "select * from macro_tasks limit 1")
	if err != nil {
		panic(err)
	}
}

func teardown() {

}
