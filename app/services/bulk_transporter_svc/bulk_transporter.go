package bulk_transporter_svc

import "gitlab.com/playment-main/support/app/models/uuid"

type ICarrier interface {
	Pick() error
	Run() error
	Drop() error
	GetId() uuid.UUID
}

type IBulkTransporter interface {
	BulkTransport([]ICarrier) (err error)
}

//--------------------------------------------------------------------------------//

//Bulk Transport logger interface
type IBulkTransportLogger interface {

	//Gets called on start of the bulk transport
	LogStartTransport()
	//Gets called on finish of bulk transport
	LogFinishTransport()

	//Gets called on individual carrier pick success
	LogPickSuccess(ICarrier)
	//Gets called on carrier pick failure
	LogPickFailure(ICarrier)
	//Gets called on carrier drop success
	LogDropSuccess(ICarrier)
	//Gets called on carrier drop failure
	LogDropFailure(ICarrier)
}

//Empty logger
type noLogger struct {
}

func (noLogger) LogStartTransport()        {}
func (noLogger) LogFinishTransport()       {}
func (noLogger) LogPickSuccess(c ICarrier) {}
func (noLogger) LogPickFailure(c ICarrier) {}
func (noLogger) LogDropSuccess(c ICarrier) {}
func (noLogger) LogDropFailure(c ICarrier) {}
