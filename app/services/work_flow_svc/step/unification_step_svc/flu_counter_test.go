package unification_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newRandomFlu() feed_line.FLU {
	return feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID: uuid.NewV4(),
		},
	}
}

func TestFluCounter_Clear(t *testing.T) {

	fc := newFluCounter()

	flu := newRandomFlu()
	flu.CopyId = 0

	fc.UpdateCount(flu)

	count := fc.GetCount(flu.ID)
	assert.Equal(t, 1, count)

	fc.Clear(flu.ID)

	count = fc.GetCount(flu.ID)
	assert.Equal(t, 0, count)
}

func TestFluCounter_GetCount(t *testing.T) {
	fc := newFluCounter()

	flu := newRandomFlu()
	flu.CopyId = 0

	flu2 := flu

	fc.UpdateCount(flu)
	fc.UpdateCount(flu2)

	assert.Equal(t, 1, fc.GetCount(flu.ID))
}

func TestFluCounter_UpdateCount(t *testing.T) {
	fc := newFluCounter()

	flu := newRandomFlu()
	flu.CopyId = 0

	flu2 := flu
	flu2.CopyId = 1

	fc.UpdateCount(flu)
	fc.UpdateCount(flu2)

	assert.Equal(t, 2, fc.GetCount(flu.ID))
}
