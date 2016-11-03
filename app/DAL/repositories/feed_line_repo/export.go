package feed_line_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
)

func New() IFluRepo {
	return &fluRepo{
		Db:       postgres.GetPostgresClient(),
		stepRepo: step_repo.New(),
	}
}

func NewInputQueue() *inputQueue {

	return &inputQueue{
		mgo: clients.GetMongoClient(),
	}
}

func NewLogger() IFluLogger {
	return &feedLineLogger{
		Db: postgres.GetPostgresClient(),
	}
}
