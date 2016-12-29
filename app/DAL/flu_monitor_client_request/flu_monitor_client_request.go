package flu_monitor_client_request

import (
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"sync"
)

// ShortHand for channel of FLUs i.e. FeedLine
type Fmcr struct {
	mq rabbitmq.MQ

	queueName string
	once      sync.Once
}

func New(name string) Fmcr {

	return Fmcr{

		mq:        rabbitmq.New(name),
		queueName: name,
	}
}

func (cr *Fmcr) Push(fmcr FMCR) {

	// Send only the models.Feedline part of the flu in bytes
	bty, _ := json.Marshal(fmcr.Flu_monitor_client_request)

	// This is async
	// TODO Think about a way to guarantee this operation also
	cr.mq.Publish(bty)

	// Just for safety: if someone forgets
	// to ConfirmReceive the flu received from a queue
	// then reconfirm it here as it will most
	// probably be a bug
	if fmcr.delivery.Acknowledger != nil {
		fmcr.ConfirmReceive()
	}

	plog.Trace("feedline", "complete push from: ", cr.queueName, "id: ", fmcr.ID.String())
}

func (cr *Fmcr) Receiver() <-chan FMCR {

	println("Feedline, subscribe request: ", cr.queueName)

	var fmcrChan chan FMCR
	var flag bool = false

	cr.once.Do(func() {

		fmcrChan = make(chan FMCR)

		go func() {

			for msg := range cr.mq.Consume() {

				fmcr := models.Flu_monitor_client_request{}
				json.Unmarshal(msg.Body, &fmcr)

				fmcrChan <- FMCR{
					Flu_monitor_client_request: fmcr,
					delivery:                   msg,
					once:                       &sync.Once{},
				}
				plog.Trace("feedline", "sent to FLU chan, name: ", cr.queueName, "id: ", fmcr.ID.String())
			}
		}()

		flag = true
	})

	if flag {
		return (<-chan FMCR)(fmcrChan)
	} else {
		panic(errors.New("Feedline already subscribed, name: " + cr.queueName))
	}

}
