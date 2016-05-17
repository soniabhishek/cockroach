package feed_line_repo

import (
	"gitlab.com/playment-main/angel/app/services/data_access_svc/clients"
	"gitlab.com/playment-main/angel/app/services/data_access_svc/repositories"
)

func New() IFluRepo {
	return NewCustom(clients.GetPostgresClient())
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
