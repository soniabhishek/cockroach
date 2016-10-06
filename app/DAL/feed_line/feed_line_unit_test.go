package feed_line

import (
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"time"
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

func TestFeedline_Load(t *testing.T) {
	fl := New("test12")

	flus := fl.Receiver()

	go func() {

		for {
			fl.Push(FLU{
				FeedLineUnit: models.FeedLineUnit{
					ID: uuid.NewV4(),
				},
			})

		}
	}()

	go func() {

		for {
			<-flus

		}
	}()

	time.Sleep(time.Duration(1) * time.Second)

	fl.amqpChan.QueueDelete("test12", false, false, false)

}
