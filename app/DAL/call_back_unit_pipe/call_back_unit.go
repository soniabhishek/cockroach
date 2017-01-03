package call_back_unit_pipe

import (
	"encoding/json"
	"errors"
	"github.com/crowdflux/angel/app/DAL/clients/rabbitmq"
	"github.com/crowdflux/angel/app/plog"
	"sync"
)

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

	bty, _ := json.Marshal(fmcr)
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

	plog.Trace("CBU", "complete push from: ", cr.queueName, "CallBackFLust: ", fmcr.FluOutputObj)
}

func (cr *CbuQ) Receiver() <-chan CBU {

	println("CBU, subscribe request: ", cr.queueName)

	var fmcrChan chan CBU
	var flag bool = false

	cr.once.Do(func() {

		fmcrChan = make(chan CBU)

		go func() {

			for msg := range cr.mq.Consume() {

				fmcr := CBU{}
				json.Unmarshal(msg.Body, &fmcr)

				fmcrChan <- CBU{
					FluOutputObj:  fmcr.FluOutputObj,
					FlusSent:      fmcr.FlusSent,
					ProjectConfig: fmcr.ProjectConfig,
					RetryLeft:     fmcr.RetryLeft,
					delivery:      msg,
					once:          &sync.Once{},
				}
				plog.Trace("CBU", "sent to FLU chan, name: ", cr.queueName, "Request: ", fmcr)
			}
		}()

		flag = true
	})

	if flag {
		return (<-chan CBU)(fmcrChan)
	} else {
		panic(errors.New("CBU already subscribed, name: " + cr.queueName))
	}

}
