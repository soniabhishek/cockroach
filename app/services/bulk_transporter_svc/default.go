package bulk_transporter_svc

type defaultTransporter struct {
	logger IBulkTransportLogger
}

func (s *defaultTransporter) BulkTransport(cs []ICarrier) error {
	s.logger.LogStartTransport()
	for _, c := range cs {
		s.transport(c)
	}
	s.logger.LogFinishTransport()
	return nil
}

func (s *defaultTransporter) transport(c ICarrier) (err error) {
	err = c.Pick()
	if err != nil {

		// We could have made the PicKFailure call in the same goroutine & let
		// the caller take care of goroutines. But the whole point of
		// logger was logging happens in a non blocking way that's why
		go s.logger.LogPickFailure(c)
		return
	}
	go s.logger.LogPickSuccess(c)

	err = c.Drop()
	if err != nil {
		go s.logger.LogDropFailure(c)
		return
	}
	go s.logger.LogDropSuccess(c)
	return
}
