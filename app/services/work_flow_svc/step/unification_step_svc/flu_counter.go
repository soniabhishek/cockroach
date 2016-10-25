package unification_step_svc

import (
	"errors"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"strconv"
	"sync"
)

type fluCounter struct {
	counter map[uuid.UUID][]feed_line.FLU

	// Read Write lock to enable concurrent
	// reads by single writes
	sync.RWMutex
}

func (f *fluCounter) GetCount(fluId uuid.UUID) int {

	f.RLock()
	defer f.RUnlock()

	return len(f.counter[fluId])
}

func (f *fluCounter) Get(fluId uuid.UUID) []feed_line.FLU {
	f.RLock()
	defer f.RUnlock()

	return f.counter[fluId]
}

func (f *fluCounter) UpdateCount(flu feed_line.FLU) {

	f.Lock()
	defer f.Unlock()

	existingFLus := f.counter[flu.ID]

	for _, eFlu := range existingFLus {

		if eFlu.StepMetaData[index] == flu.StepMetaData[index] {

			index := eFlu.StepMetaData[index].(int)
			indexStr := strconv.Itoa(index)
			plog.Error("FLU Counter", errors.New("Already updated counter for flu_id : "+flu.ID.String()+" index : "+indexStr))
			return
		}
	}

	f.counter[flu.ID] = append(existingFLus, flu)
}

func (f *fluCounter) Clear(fluId uuid.UUID) {

	f.Lock()
	defer f.Unlock()

	flus, ok := f.counter[fluId]
	if !ok {
		return
	}

	for _, flu := range flus {
		flu.ConfirmReceive()
	}

	delete(f.counter, fluId)
}

func newFluCounter() fluCounter {
	return fluCounter{make(map[uuid.UUID][]feed_line.FLU), sync.RWMutex{}}
}
