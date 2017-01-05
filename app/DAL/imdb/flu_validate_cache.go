package imdb

import (
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/models"
)

type FluValidateCache struct {
	imdb IAngelImdb
}

func (f *FluValidateCache) Set(key uuid.UUID, val []models.FLUValidator) {
	f.imdb.Set(key.String(), val)
}
func (f *FluValidateCache) Get(key uuid.UUID) (val []models.FLUValidator, err error) {
	obj, err := f.imdb.Get(key.String())
	if err != nil {
		return
	}
	return obj.([]models.FLUValidator), nil
}

func (f *FluValidateCache) Remove(key uuid.UUID) error {
	return f.imdb.Remove(key.String())
}

func (f *FluValidateCache) SafeSet(key uuid.UUID, val []models.FLUValidator) (res []models.FLUValidator, err error) {
	obj, err := f.imdb.SafeSet(key.String(), val)
	if err != nil {
		return
	}
	return obj.([]models.FLUValidator), nil
}

func (f *FluValidateCache) ClearAll() {
	f.imdb.ClearAll()
}

