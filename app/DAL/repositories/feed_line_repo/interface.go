package feed_line_repo

import (
	"errors"
	"sync"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IFluRepo interface {
	GetById(Id uuid.UUID) (models.FeedLineUnit, error)
	Save(feedLineUnit models.FeedLineUnit)
	BulkInsert(flus []models.FeedLineUnit) error
	BulkUpdate(flus []models.FeedLineUnit) error
	Add(feedLineUnit models.FeedLineUnit) error
	Update(feedLineUnit models.FeedLineUnit) error
	GetByStepId(StepId uuid.UUID) ([]models.FeedLineUnit, error)
	BulkFluBuildUpdate(flus []models.FeedLineUnit) error
	BulkFluBuildUpdateByStepType(flus []models.FeedLineUnit, stepType step_type.StepType) (updatedFlus []models.FeedLineUnit, err error)
}

type IFluLogger interface {
	Log([]models.FeedLineLog) error
}

//--------------------------------------------------------------------------------//

type inMemFluRepo struct {
	sync.RWMutex
	flus map[uuid.UUID]models.FeedLineUnit
}

var _ IFluRepo = &inMemFluRepo{}

func (i *inMemFluRepo) GetById(Id uuid.UUID) (models.FeedLineUnit, error) {
	i.RLock()
	defer i.RUnlock()

	flu, ok := i.flus[Id]
	if !ok {
		return models.FeedLineUnit{}, ErrFLUNotFound
	}
	return flu, nil
}

func (i *inMemFluRepo) Save(flu models.FeedLineUnit) {
	i.Lock()
	defer i.Unlock()

	i.flus[flu.ID] = flu
}

func (i *inMemFluRepo) BulkInsert(flus []models.FeedLineUnit) error {
	i.Lock()
	defer i.Unlock()

	for _, flu := range flus {
		i.flus[flu.ID] = flu
	}
	return nil
}

func (i *inMemFluRepo) BulkUpdate(flusToUpdate []models.FeedLineUnit) error {
	i.Lock()
	defer i.Unlock()

	for _, flu := range flusToUpdate {
		_, ok := i.flus[flu.ID]
		if !ok {
			return errors.New(flu.ID.String() + " not present")
		} else {
			i.flus[flu.ID] = flu
		}
	}
	return nil
}

func (i *inMemFluRepo) Add(flu models.FeedLineUnit) error {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.flus[flu.ID]; ok {
		return errors.New("already present")
	}
	i.flus[flu.ID] = flu
	return nil
}

func (i *inMemFluRepo) Update(flu models.FeedLineUnit) error {
	i.Lock()
	defer i.Unlock()

	if _, ok := i.flus[flu.ID]; !ok {
		return errors.New("not present")
	}

	i.flus[flu.ID] = flu
	return nil
}

func (i *inMemFluRepo) GetByStepId(stepId uuid.UUID) (stepFlus []models.FeedLineUnit, err error) {
	i.RLock()
	defer i.RUnlock()

	for _, flu := range i.flus {
		stepFlus = append(stepFlus, flu)
	}
	return stepFlus, nil
}

func (e *inMemFluRepo) BulkFluBuildUpdate(flus []models.FeedLineUnit) error {

	return nil
}

func (e *inMemFluRepo) BulkFluBuildUpdateByStepType(flus []models.FeedLineUnit, stepType step_type.StepType) (updatedFlus []models.FeedLineUnit, err error) {
	return
}

func Mock() IFluRepo {
	return &inMemFluRepo{
		flus: make(map[uuid.UUID]models.FeedLineUnit),
	}
}
