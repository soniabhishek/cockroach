package bulk_transporter_svc

func NewWithNoLogger() IBulkTransporter {
	return &defaultTransporter{
		logger: noLogger{},
	}
}

func New(logger IBulkTransportLogger) IBulkTransporter {
	return &defaultTransporter{
		logger: logger,
	}
}
