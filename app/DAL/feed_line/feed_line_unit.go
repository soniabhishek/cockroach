package feed_line

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/streadway/amqp"
	"sync"
)

//--------------------------------------------------------------------------------//

type FLU struct {
	models.FeedLineUnit

	delivery amqp.Delivery

	once sync.Once
}

func (flu *FLU) ConfirmReceive() {

	flu.once.Do(func() {
		err := flu.delivery.Ack(false)
		if err != nil {
			plog.Error("FLU", err, "error while ack", "fluId: "+flu.FeedLineUnit.ID.String())
			panic(err)
		}
	})
}

func (flu *FLU) Redelivered() bool {
	return flu.delivery.Redelivered
}
