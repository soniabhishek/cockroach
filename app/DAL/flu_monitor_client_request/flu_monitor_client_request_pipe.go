package flu_monitor_client_request

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/streadway/amqp"
	"sync"
)

//--------------------------------------------------------------------------------//

type FMCR struct {
	models.Flu_monitor_client_request

	delivery amqp.Delivery

	once *sync.Once
}

func (fmcr *FMCR) ConfirmReceive() {

	defer func() {
		// to handle if flu.once is nil
		recover()
	}()

	fmcr.once.Do(func() {

		err := fmcr.delivery.Ack(false)
		if err != nil {
			plog.Error("FMCR", err, "error while ack", "fmcrId: "+fmcr.Flu_monitor_client_request.ID.String())
			panic(err)
		}
	})
}

func (fmcr *FMCR) Redelivered() bool {
	return fmcr.delivery.Redelivered
}

/*
func (flu FMCR) Copy() FMCR {
	flu.Build = flu.Build.Copy()

	flu.delivery = amqp.Delivery{}
	flu.once = &sync.Once{}
	return flu
}
*/
