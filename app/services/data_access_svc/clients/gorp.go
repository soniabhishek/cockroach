package clients

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
	"gitlab.com/playment-main/support/app/config"
	"gitlab.com/playment-main/support/app/models"
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

//--------------------------------------------------------------------------------//

//This wrapper postgres_db is to have one place where we can write all the modifications on errors thrown by gorp & pq

type postgres_db struct {
	gorpDbMap *gorp.DbMap
}

func (pg *postgres_db) Insert(list ...interface{}) (err error) {
	err = pg.gorpDbMap.Insert(list...)
	return
}
func (pg *postgres_db) Update(list ...interface{}) (d int64, err error) {
	return pg.gorpDbMap.Update(list...)
}
func (pg *postgres_db) Delete(list ...interface{}) (int64, error) {
	return pg.gorpDbMap.Delete(list...)
}
func (pg *postgres_db) Exec(query string, args ...interface{}) (sql.Result, error) {
	return pg.gorpDbMap.Exec(query, args...)
}
func (pg *postgres_db) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return pg.gorpDbMap.Select(i, query, args...)
}
func (pg *postgres_db) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return pg.gorpDbMap.SelectOne(holder, query, args...)
}

func GetPostgresClient() *postgres_db {
	return &postgres_db{gorpDbMap}
}
