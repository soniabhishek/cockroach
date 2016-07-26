package clients

import (
	"github.com/crowdflux/angel/app/config"
	"gopkg.in/mgo.v2"
)

var mongo_db *mgo.Database

func init() {
	mongo_db = initMongoDb()
}

func initMongoDb() *mgo.Database {

	session, err := mgo.Dial(config.MONGO_HOST.Get())
	if err != nil {
		panic(err)
	}
	//session.SetMode(mgo.Monotonic, true)

	db := session.DB(config.MONGO_DB_NAME.Get())

	return db
}

func GetMongoClient() *mgo.Database {
	return mongo_db
}
