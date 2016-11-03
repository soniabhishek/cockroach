package feed_line

import (
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"sync"
)

// ShortHand for channel of FLUs i.e. FeedLine
type Fl struct {
	mq rabbitmq.MQ

	queueName string
	once      sync.Once
}

func New(name string) Fl {

	return Fl{

		mq:        rabbitmq.New(name),
		queueName: name,
	}
}

func (fl *Fl) Push(flu FLU) {

	// Send only the models.Feedline part of the flu in bytes
	bty, _ := json.Marshal(flu.FeedLineUnit)

	// This is async
	// TODO Think about a way to guarantee this operation also
	fl.mq.Publish(bty)

	// Just for safety: if someone forgets
	// to ConfirmReceive the flu received from a queue
	// then reconfirm it here as it will most
	// probably be a bug
	if flu.delivery.Acknowledger != nil {
		flu.ConfirmReceive()
	}

	plog.Trace("feedline", "complete push from: ", fl.queueName, "id: ", flu.ID.String())
}

func (fl *Fl) Receiver() <-chan FLU {

	println("Feedline, subscribe request: ", fl.queueName)

	var fluChan chan FLU
	var flag bool = false

	fl.once.Do(func() {

		fluChan = make(chan FLU)

		go func() {

			for msg := range fl.mq.Consume() {

				flu := models.FeedLineUnit{}
				json.Unmarshal(msg.Body, &flu)

				fluChan <- FLU{
					FeedLineUnit: flu,
					delivery:     msg,
					once:         &sync.Once{},
				}
				plog.Trace("feedline", "sent to FLU chan, name: ", fl.queueName, "id: ", flu.ID.String())
			}
		}()

		flag = true
	})

	if flag {
		return (<-chan FLU)(fluChan)
	} else {
		panic(errors.New("Feedline already subscribed, name: " + fl.queueName))
	}

}
