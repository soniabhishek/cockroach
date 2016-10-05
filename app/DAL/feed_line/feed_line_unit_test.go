package feed_line

import (
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	fl := New("test1")

	fluId := uuid.NewV4()

	fl.Push(FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID: fluId,
		},
	})

	flu := <-fl.Receiver()

	flu.ConfirmReceive()

	fl.amqpChan.QueueDelete("test1", false, false, false)

	assert.EqualValues(t, fluId, flu.ID)

}
