package plog_logger

type ILogger interface {
	Write(message ...interface{})
}
