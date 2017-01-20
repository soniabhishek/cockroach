package imdb

import "github.com/crowdflux/angel/app/models"

type fluUploadCache struct {
	imdb IAngelImdb
}

func (f *fluUploadCache) Set(key string, fUS models.FluUploadStats) {
	f.imdb.Set(key, fUS)
}
func (f *fluUploadCache) Get(key string) (fUS models.FluUploadStats, err error) {
	obj, err := f.imdb.Get(key)
	if err != nil {
		return
	}
	return obj.(models.FluUploadStats), nil
}

func (f *fluUploadCache) Remove(key string) error {
	return f.imdb.Remove(key)
}
func (f *fluUploadCache) SafeSet(key string, value models.FluUploadStats) (flu models.FluUploadStats, err error) {
	obj, err := f.imdb.SafeSet(key, value)
	if err != nil {
		return
	}
	return obj.(models.FluUploadStats), nil
}

func (f *fluUploadCache) ClearAll() {
	f.imdb.ClearAll()
}
