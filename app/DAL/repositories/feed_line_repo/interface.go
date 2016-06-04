package feed_line_repo

import (
	"errors"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"sync"
)

type IFluRepo interface {
	GetById(Id uuid.UUID) (models.FeedLineUnit, error)
	Save(feedLineUnit models.FeedLineUnit)
	BulkInsert(flus []models.FeedLineUnit) error
	Add(feedLineUnit models.FeedLineUnit) error
	Update(feedLineUnit models.FeedLineUnit) error
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

func Mock() IFluRepo {
	return &inMemFluRepo{
		flus: make(map[uuid.UUID]models.FeedLineUnit),
	}
}
