package feed_line_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
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
