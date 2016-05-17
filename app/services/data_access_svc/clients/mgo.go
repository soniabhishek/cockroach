package clients

import (
	"sync"

	"gitlab.com/playment-main/angel/app/config"
	"gopkg.in/mgo.v2"
)

var mongo_db *mgo.Database
var onceMgo sync.Once

func init() {
	onceMgo.Do(func() {
		mongo_db = initMongoDb()
	})
}

func initMongoDb() *mgo.Database {

	session, err := mgo.Dial(config.GetVal(config.MONGO_HOST))
	if err != nil {
		panic(err)
	}
	//session.SetMode(mgo.Monotonic, true)

	//Config package not working as expected
	//again requiring setEnvironment
	config.SetEnvironment(config.Development)

	db := session.DB(config.GetVal(config.MONGO_DB_NAME))

	return db
}

func GetMongoClient() *mgo.Database {
	return mongo_db
}
