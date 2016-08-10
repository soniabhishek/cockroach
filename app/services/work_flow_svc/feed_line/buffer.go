package feed_line

import (
	"github.com/crowdflux/angel/app/models/uuid"
	"sync"
)

type Bf struct {
	mtx    sync.RWMutex
	fluMap map[uuid.UUID]FLU
}

func NewBuffer() Bf {
	return Bf{fluMap: make(map[uuid.UUID]FLU)}
}

// RLock is read lock i.e. either multiple reads
// or single write can happen at a time
func (b *Bf) Get(id uuid.UUID) (FLU, bool) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	flu, ok := b.fluMap[id]
	return flu, ok
}

// Write lock part of the read write lock
func (b *Bf) Save(flu FLU) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.fluMap[flu.ID] = flu
}

func (b *Bf) GetAll() map[uuid.UUID]FLU {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	return b.fluMap
}
