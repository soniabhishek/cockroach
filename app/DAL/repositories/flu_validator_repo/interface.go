package flu_validator_repo

import (
	"errors"
	"sync"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IFluValidatorRepo interface {
	GetValidatorsForProject(projectId uuid.UUID, tag string) ([]models.FLUValidator, error)
	Save(*models.FLUValidator) error
}

//--------------------------------------------------------------------------------//

type inMemFluValidatorRepo struct {
	sync.RWMutex
	fluVs map[uuid.UUID]models.FLUValidator
}

var _ IFluValidatorRepo = &inMemFluValidatorRepo{}

func (i *inMemFluValidatorRepo) GetValidatorsForProject(projectId uuid.UUID, tag string) ([]models.FLUValidator, error) {
	i.RLock()
	defer i.RUnlock()

	projVs := []models.FLUValidator{}

	for _, v := range i.fluVs {

		if v.ProjectId == projectId && v.Tag == tag {
			projVs = append(projVs)
		}
	}

	if len(projVs) == 0 {
		return projVs, errors.New("not found")
	}

	return projVs, nil
}

func (i *inMemFluValidatorRepo) GetById(Id uuid.UUID) (models.FLUValidator, error) {
	i.RLock()
	defer i.RUnlock()

	v, ok := i.fluVs[Id]
	if !ok {
		return models.FLUValidator{}, errors.New("not found")
	}
	return v, nil
}

func (i *inMemFluValidatorRepo) Save(v *models.FLUValidator) error {
	i.Lock()
	defer i.Unlock()

	i.fluVs[v.ID] = *v
	return nil
}

func (i *inMemFluValidatorRepo) BulkInsert(fluVs []models.FLUValidator) error {
	i.Lock()
	defer i.Unlock()

	for _, v := range fluVs {
		i.fluVs[v.ID] = v
	}
	return nil
}

func (i *inMemFluValidatorRepo) Add(v models.FLUValidator) error {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.fluVs[v.ID]; ok {
		return errors.New("already present")
	}
	i.fluVs[v.ID] = v
	return nil
}

func (i *inMemFluValidatorRepo) Update(v models.FLUValidator) error {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.fluVs[v.ID]; !ok {
		return errors.New("not present")
	}

	i.fluVs[v.ID] = v
	return nil
}

func Mock() IFluValidatorRepo {
	return &inMemFluValidatorRepo{
		fluVs: make(map[uuid.UUID]models.FLUValidator),
	}
}
