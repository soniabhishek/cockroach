package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/config"
)

var sqlxDB *sqlx.DB

func init() {
	sqlxDB = initSQLxClient()
}

func initSQLxClient() *sqlx.DB {

	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	dbName := config.Get(config.PG_DATABASE_NAME)

	db := sqlx.MustConnect("postgres", "dbname="+dbName+" user=postgres password=postgres host=localhost sslmode=disable")
	return db
}

//--------------------------------------------------------------------------------//

type sqlXWrap struct {
}

func (sqlXWrap) SelectCustom(tables []interface{}, query string, args ...interface{}) {
	//sqlxDB.QueryRow()
}

func (sqlXWrap) SelectOneCustom(tables []interface{}, query string, args ...interface{}) {
	//sqlxDB.Query()
}
