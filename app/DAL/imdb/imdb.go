package imdb

import (
	"github.com/pkg/errors"
	"sync"
)

type imdb struct {
	db map[string]interface{}
	rw sync.RWMutex
}

func (imdb *imdb) Get(key string) (interface{}, error) {
	imdb.rw.RLock()
	defer imdb.rw.RUnlock()
	val, ok := imdb.db[key]
	if !ok {
		return nil, errors.New("NO Such Key Exist")
	} else {
		return val, nil
	}
}

func (imdb *imdb) Set(key string, value interface{}) {
	imdb.rw.Lock()
	defer imdb.rw.Unlock()
	imdb.db[key] = value

}

func (imdb *imdb) Remove(key string) error {
	imdb.rw.Lock()
	defer imdb.rw.Unlock()
	_, ok := imdb.db[key]
	if ok {
		delete(imdb.db, key)
		return nil
	} else {
		return errors.New("No Such Key Exist")
	}
}

//Returns error if key already exist in cache.
func (imdb *imdb) SafeSet(key string, value interface{}) (interface{}, error) {
	imdb.rw.Lock()
	defer imdb.rw.Unlock()
	val, ok := imdb.db[key]
	if !ok {
		imdb.db[key] = value
		return nil, nil
	} else {
		return val, errors.New("Key Already Exist")
	}
}

func (imdb *imdb) ClearAll() {
	imdb.rw.Lock()
	defer imdb.rw.Unlock()
	imdb.db = make(map[string]interface{})
}
