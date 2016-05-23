package postgres

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"gopkg.in/gorp.v1"
)

var gorpDbMap *gorp.DbMap
var onceGorp sync.Once

func init() {
	onceGorp.Do(func() {
		gorpDbMap = initGorpClient()
	})
}

func initGorpClient() *gorp.DbMap {

	//Config package not working as expected
	//again requiring setEnvironment
	config.SetEnvironment(config.Development)

	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	dbName := config.GetVal(config.DB_DATABASE_NAME)

	db, err := sql.Open("postgres", "dbname="+dbName+" user=postgres password=postgres host=localhost sslmode=disable")
	if err != nil {
		panic("Main db connection failed")
	}

	db.SetMaxIdleConns(40)
	db.SetMaxOpenConns(160)

	// construct a gorp DbMap
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	addTableInfo(dbMap)

	return dbMap
}

func addTableInfo(dbMap *gorp.DbMap) {
	dbMap.AddTableWithName(models.FeedLineUnit{}, "feed_line").SetKeys(false, "id")
	dbMap.AddTableWithName(models.FLUValidator{}, "input_flu_validator").SetKeys(false, "id")
	dbMap.AddTableWithName(models.MacroTask{}, "macro_tasks").SetKeys(false, "id")
	dbMap.AddTableWithName(models.Project{}, "projects").SetKeys(false, "id")
}

func GetGorpClient() *gorp.DbMap {
	return gorpDbMap
}
