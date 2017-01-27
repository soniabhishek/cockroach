package imdb

import (
	"github.com/crowdflux/angel/app/models"
	"sync"
)

type cmap struct {
	cmap map[interface{}]interface{}
	rw   sync.RWMutex
}

func (c *cmap) Get(key interface{}) (val interface{}, ok bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	val, ok = c.cmap[key]
	return
}

func (c *cmap) Set(key, value interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.cmap[key] = value
}

func (c *cmap) Delete(key interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	delete(c.cmap, key)
}

func (c *cmap) Reset() {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.cmap = make(map[interface{}]interface{})
}

func (c *cmap) Iter() <-chan models.Tuple {
	ch := make(chan models.Tuple)
	c.rw.RLock()
	for key, value := range c.cmap {
		ch <- models.Tuple{key, value}
	}
	c.rw.RUnlock()
	close(ch)
	return ch
}
