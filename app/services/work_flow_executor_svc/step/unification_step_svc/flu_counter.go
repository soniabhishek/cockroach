package unification_step_svc

import (
	"errors"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
	"sync"
)

type fluStepGroup struct {
	MasterFluId uuid.UUID
	StepId      uuid.UUID
}

type fluCounter struct {
	counter map[fluStepGroup][]feed_line.FLU

	// Read Write lock to enable concurrent
	// reads and single writes
	sync.RWMutex
}

func (f *fluCounter) GetCount(flu feed_line.FLU) int {

	f.RLock()
	defer f.RUnlock()

	return len(f.counter[getFluStepGroup(flu)])
}

func (f *fluCounter) Get(flu feed_line.FLU) []feed_line.FLU {
	f.RLock()
	defer f.RUnlock()

	return f.counter[getFluStepGroup(flu)]
}

func (f *fluCounter) UpdateCount(flu feed_line.FLU) {

	f.Lock()
	defer f.Unlock()

	existingFLus := f.counter[getFluStepGroup(flu)]

	for _, eFlu := range existingFLus {

		if eFlu.ID == flu.ID {

			plog.Error("FLU Counter", errors.New("Already updated counter for flu_id : "+flu.ID.String()), plog.MessageWithParam(log_tags.MASTER_FLU_ID, flu.MasterId.String()))
			return
		}
	}

	f.counter[getFluStepGroup(flu)] = append(existingFLus, flu)
}

func (f *fluCounter) Clear(flu feed_line.FLU) {

	f.Lock()
	defer f.Unlock()

	flus, ok := f.counter[getFluStepGroup(flu)]
	if !ok {
		return
	}

	for _, flu := range flus {
		flu.ConfirmReceive()
	}

	delete(f.counter, getFluStepGroup(flu))
}

func newFluCounter() fluCounter {
	return fluCounter{make(map[fluStepGroup][]feed_line.FLU), sync.RWMutex{}}
}

func getFluStepGroup(flu feed_line.FLU) fluStepGroup {
	return fluStepGroup{flu.MasterId, flu.StepId}
}
