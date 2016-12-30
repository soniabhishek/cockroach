package call_back_unit_pipe

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/streadway/amqp"
	"net/http"
	"sync"
)

//--------------------------------------------------------------------------------//

type CBU struct {
	FluOutputObj []models.FluOutputStruct

	FlusSent map[uuid.UUID]feed_line.FLU

	ProjectConfig models.ProjectConfiguration

	RetryCount int

	delivery amqp.Delivery

	once *sync.Once
}

func (fmcr *CBU) ConfirmReceive() {

	defer func() {
		// to handle if flu.once is nil
		recover()
	}()

	fmcr.once.Do(func() {

		err := fmcr.delivery.Ack(false)
		if err != nil {
			plog.Error("CBU", err, "error while ack", "RequestOject: ", fmcr.FluOutputObj)
			panic(err)
		}
	})
}

func (fmcr *CBU) Redelivered() bool {
	return fmcr.delivery.Redelivered
}

/*
func (flu CBU) Copy() CBU {
	flu.Build = flu.Build.Copy()

	flu.delivery = amqp.Delivery{}
	flu.once = &sync.Once{}
	return flu
}
*/
