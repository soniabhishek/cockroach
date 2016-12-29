package flu_monitor_client_request

import (
	"github.com/crowdflux/angel/app/models/uuid"
	"sync"
)

type Bf struct {
	mtx     sync.RWMutex
	fmcrMap map[uuid.UUID]FMCR
}

func NewBuffer() Bf {
	return Bf{fmcrMap: make(map[uuid.UUID]FMCR)}
}

// RLock is read lock i.e. either multiple reads
// or single write can happen at a time
func (b *Bf) Get(id uuid.UUID) (FMCR, bool) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	fmcr, ok := b.fmcrMap[id]
	return fmcr, ok
}

// Write lock part of the read write lock
func (b *Bf) Save(fmcr FMCR) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.fmcrMap[fmcr.ID] = fmcr
}

func (b *Bf) GetAll() map[uuid.UUID]FMCR {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	return b.fmcrMap
}
