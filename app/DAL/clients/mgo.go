package clients

import (
	"gitlab.com/playment-main/angel/app/config"
	"gopkg.in/mgo.v2"
)

var mongo_db *mgo.Database

func init() {
	mongo_db = initMongoDb()
}

func initMongoDb() *mgo.Database {

	session, err := mgo.Dial(config.Get(config.MONGO_HOST))
	if err != nil {
		panic(err)
	}
	//session.SetMode(mgo.Monotonic, true)

	db := session.DB(config.Get(config.MONGO_DB_NAME))

	return db
}

func GetMongoClient() *mgo.Database {
	return mongo_db
}
