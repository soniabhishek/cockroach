package imdb

import (
	"github.com/pkg/errors"
	"sync"
)

type imdb struct {
	db cmap
	rw sync.RWMutex
}

func (imdb *imdb) Get(key string) (interface{}, error) {
	val, ok := imdb.db.Get(key)
	if !ok {
		return nil, errors.New("NO Such Key Exist")
	} else {
		return val, nil
	}
}

func (imdb *imdb) Set(key string, value interface{}) {
	imdb.db.Set(key, value)
}

func (imdb *imdb) Remove(key string) error {
	_, ok := imdb.db.Get(key)
	if ok {
		imdb.db.Delete(key)
		return nil
	} else {
		return errors.New("No Such Key Exist")
	}
}

//Returns error if key already exist in cache.
func (imdb *imdb) SafeSet(key string, value interface{}) (interface{}, error) {
	val, ok := imdb.db.Get(key)
	if !ok {
		imdb.db.Set(key, value)
		return nil, nil
	} else {
		return val, errors.New("Key Already Exist")
	}
}

func (imdb *imdb) ClearAll() {
	imdb.db.Reset()
}
