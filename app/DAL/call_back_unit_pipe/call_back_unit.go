package call_back_unit_pipe

import (
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"sync"
)

// ShortHand for channel of FLUs i.e. FeedLine
type CbuQ struct {
	mq        rabbitmq.MQ
	queueName string
	once      sync.Once
}

func New(name string) CbuQ {

	return CbuQ{

		mq:        rabbitmq.New(name),
		queueName: name,
	}
}

func (cr *CbuQ) Push(fmcr CBU) {

	// Send only the models.Feedline part of the flu in bytes
	bty, _ := json.Marshal(fmcr.FluOutputObj)

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

	plog.Trace("feedline", "complete push from: ", cr.queueName, "CallBackFLust: ", fmcr.FluOutputObj)
}

func (cr *CbuQ) Receiver() <-chan CBU {

	println("Feedline, subscribe request: ", cr.queueName)

	var fmcrChan chan CBU
	var flag bool = false

	cr.once.Do(func() {

		fmcrChan = make(chan CBU)

		go func() {

			for msg := range cr.mq.Consume() {

				fmcr := []models.FluOutputStruct{}
				json.Unmarshal(msg.Body, &fmcr)

				fmcrChan <- CBU{
					FluOutputObj: fmcr,
					delivery:     msg,
					once:         &sync.Once{},
				}
				plog.Trace("feedline", "sent to FLU chan, name: ", cr.queueName, "Request: ", fmcr)
			}
		}()

		flag = true
	})

	if flag {
		return (<-chan CBU)(fmcrChan)
	} else {
		panic(errors.New("Feedline already subscribed, name: " + cr.queueName))
	}

}
