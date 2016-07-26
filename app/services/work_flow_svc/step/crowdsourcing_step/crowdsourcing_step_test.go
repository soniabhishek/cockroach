package crowdsourcing_step

import (
	"testing"

	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
	"github.com/stretchr/testify/assert"
	//"time"
	"time"
)

//--------------------------------------------------------------------------------//

type fakeFluPusher struct {
}

func (fakeFluPusher) PushFLU(models.FeedLineUnit) (bool, error) {
	return true, nil
}

var fluId = uuid.NewV4()

var flu = feed_line.FLU{
	FeedLineUnit: models.FeedLineUnit{
		ID:          fluId,
		ReferenceId: "PayFlip123",
		Tag:         "Brand",
		Data: models.JsonFake{
			"brand":  "Sony",
			"image1": "http://sxomeimaghe.com/some.jpeg",
		},
		Build: models.JsonFake{},
	},
}

func Test(t *testing.T) {

	fluRepo := feed_line_repo.Mock()

	fluRepo.Save(flu.FeedLineUnit)

	cs := crowdSourcingStep{
		Step:      step.New(),
		fluRepo:   feed_line_repo.New(),
		fluClient: fakeFluPusher{},
	}

	cs.start()

	cs.InQ <- flu

	// Giving it time to finish adding to buffer
	// as its happening in another goroutine
	time.Sleep(time.Duration(100) * time.Millisecond)

	flu.Build["new_prop"] = 123

	ok := cs.finishFlu(flu)

	assert.True(t, ok)

	var fluNew feed_line.FLU
	select {
	case fluNew = <-cs.OutQ:
		assert.EqualValues(t, flu.ID, fluNew.ID)
		assert.EqualValues(t, flu.Build["new_prop"], 123)
	default:
		assert.FailNow(t, "nothing came out of crowdsourcing queue")
	}

}

func TestInvalidFlu(t *testing.T) {

	fluRepo := feed_line_repo.Mock()

	fluRepo.Save(flu.FeedLineUnit)

	cs := crowdSourcingStep{
		Step:      step.New(),
		fluRepo:   feed_line_repo.New(),
		fluClient: fakeFluPusher{},
	}

	cs.start()

	cs.InQ <- flu

	inValidFlu := flu
	inValidFlu.ID = uuid.NewV4()

	ok := cs.finishFlu(inValidFlu)

	assert.False(t, ok)
}
