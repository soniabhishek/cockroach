package flu_validator_repo

import (
	"time"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/DAL/repositories/queries"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
)

const fluValidatorTable string = "input_flu_validator"

type fluValidatorRepo struct {
	db repositories.IDatabase
}

var _ IFluValidatorRepo = &fluValidatorRepo{}

func (f *fluValidatorRepo) GetValidatorsForProject(projectId uuid.UUID, tag string) (validators []models.FLUValidator, err error) {
	_, err = f.db.Select(&validators, "select * from input_flu_validator where project_id = $1 and tag = $2", projectId, tag)
	return
}

func (f *fluValidatorRepo) Save(fv *models.FLUValidator) error {
	existingFv, err := f.get(fv.ID)

	var exists bool

	if err == nil {
		exists = true
	}

	if exists {

		//It requires to copy the existing createdAt otherwise it will update it as nil
		fv.CreatedAt = existingFv.CreatedAt
		fv.UpdatedAt = pq.NullTime{time.Now(), true}
		_, err = f.db.Update(fv)
	} else {

		fv.CreatedAt = pq.NullTime{time.Now(), true}
		if fv.ID == uuid.Nil {
			fv.ID = uuid.NewV4()
		}
		err = f.db.Insert(fv)
	}
	return err
}

func (f *fluValidatorRepo) get(id uuid.UUID) (v models.FLUValidator, err error) {

	err = f.db.SelectOne(&v, queries.SelectById(fluValidatorTable), id)
	return
}

func (f *fluValidatorRepo) update(fv models.FLUValidator) error {

	_, err := f.db.Update(&fv)
	return err
}

func (f *fluValidatorRepo) insertMany(validators []models.FLUValidator) error {

	err := f.db.Insert(toInterfaceArray(validators)...)
	return err
}

func (f *fluValidatorRepo) deleteMany(validators []models.FLUValidator) error {

	_, err := f.db.Delete(toInterfaceArray(validators)...)
	return err
}

func toInterfaceArray(fvs []models.FLUValidator) []interface{} {
	var validatorsP []interface{} = make([]interface{}, len(fvs))
	for i, _ := range fvs {
		validatorsP[i] = &fvs[i]
	}
	return validatorsP
}
