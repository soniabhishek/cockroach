package feed_line_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/DAL/repositories"
)

func New() IFluRepo {
	return NewCustom(postgres.GetPostgresClient())
}

func NewCustom(dbInterface repositories.IDatabase) IFluRepo {
	return &fluRepo{
		Db: dbInterface,
	}
}

func NewInputQueue() *inputQueue {

	return &inputQueue{
		mgo: clients.GetMongoClient(),
	}
}
