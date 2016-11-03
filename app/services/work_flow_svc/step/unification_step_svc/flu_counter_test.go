package unification_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newRandomFlu() feed_line.FLU {
	id := uuid.NewV4()
	return feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:       id,
			Build:    models.JsonF{},
			StepId:   uuid.NewV4(),
			IsMaster: true,
			MasterId: id,
		},
	}
}

func TestFluCounter_Clear(t *testing.T) {

	fc := newFluCounter()

	flu := newRandomFlu()

	fc.UpdateCount(flu)

	count := fc.GetCount(flu)
	assert.Equal(t, 1, count)

	fc.Clear(flu)

	count = fc.GetCount(flu)
	assert.Equal(t, 0, count)
}

func TestFluCounter_GetCount(t *testing.T) {
	fc := newFluCounter()

	flu := newRandomFlu()

	flu2 := flu.Copy()

	fc.UpdateCount(flu)
	fc.UpdateCount(flu2)

	assert.Equal(t, 1, fc.GetCount(flu))
}

func TestFluCounter_UpdateCount(t *testing.T) {
	fc := newFluCounter()

	flu := newRandomFlu()

	flu2 := flu.Copy()
	flu.ID = uuid.NewV4()
	flu.IsMaster = false
	flu2.MasterId = flu.MasterId

	fc.UpdateCount(flu)
	fc.UpdateCount(flu2)

	assert.Equal(t, 2, fc.GetCount(flu))
}
