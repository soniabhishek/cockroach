package postgres

import (
	"database/sql"

	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

var gorpDbMap *gorp.DbMap

func init() {
	gorpDbMap = initGorpClient()
}

func initGorpClient() *gorp.DbMap {

	dbName := config.PG_DATABASE_NAME.Get()
	username := config.PG_USERNAME.Get()
	host := config.PG_HOST.Get()
	password := config.PG_PASSWORD.Get()

	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", `dbname=`+dbName+` user=`+username+` password=`+password+` host=`+host+` sslmode=disable`)
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
	dbMap.AddTableWithName(models.Question{}, "questions").SetKeys(false, "id")
	dbMap.AddTableWithName(models.Step{}, "step").SetKeys(false, "id")
	dbMap.AddTableWithName(models.Route{}, "routes").SetKeys(false, "id")
	dbMap.AddTableWithName(models.LogicGate{}, "logic_gate").SetKeys(false, "id")
	dbMap.AddTableWithName(models.FeedLineLog{}, "feed_line_log").SetKeys(true, "id")
	dbMap.AddTableWithName(models.User{}, "users").SetKeys(false, "id")
	dbMap.AddTableWithName(models.Client{}, "clients").SetKeys(false, "id")
	dbMap.AddTableWithName(models.ProjectConfiguration{}, "project_configuration").SetKeys(false, "project_id")
	dbMap.AddTableWithName(models.WorkFlow{}, "work_flow").SetKeys(false, "id")
}

func GetGorpClient() *gorp.DbMap {
	return gorpDbMap
}
