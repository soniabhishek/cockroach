package imdb

import "sync"

type cmap struct {
	cmap map[interface{}]interface{}
	rw   sync.RWMutex
}

func (c *cmap) get(key interface{}) (val interface{}, ok bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	val, ok = c.cmap[key]
	return
}

func (c *cmap) set(key, value interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.cmap[key] = value
}

func (c *cmap) delete(key interface{}) {
	c.rw.Lock()
	defer c.rw.Unlock()
	delete(c.cmap, key)
}

func (c *cmap) reset() {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.cmap = make(map[interface{}]interface{})
}
