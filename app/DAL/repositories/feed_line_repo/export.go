package feed_line_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
)

func New() IFluRepo {
	return &fluRepo{
		Db: postgres.GetPostgresClient(),
	}
}

func NewInputQueue() *inputQueue {

	return &inputQueue{
		mgo: clients.GetMongoClient(),
	}
}

func newLogger() IFluLogger {
	return &feedLineLogger{
		Db: postgres.GetPostgresClient(),
	}
}

var StdLogger = newLogger()
