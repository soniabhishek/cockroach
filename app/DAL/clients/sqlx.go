package clients

import (
	"sync"

	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gitlab.com/playment-main/angel/app/config"
)

var onceSQLx sync.Once
var sqlxDB *sqlx.DB

func init() {
	onceSQLx.Do(func() {
		sqlxDB = initSQLxClient()
	})
}

func initSQLxClient() *sqlx.DB {

	//Config package not working as expected
	//again requiring setEnvironment
	config.SetEnvironment(config.Development)

	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	dbName := config.GetVal(config.DB_DATABASE_NAME)

	db := sqlx.MustConnect("postgres", "dbname="+dbName+" user=postgres password=postgres host=localhost sslmode=disable")
	return db
}

//--------------------------------------------------------------------------------//

type sqlXWrap struct {
}

func (*sqlXWrap) Insert(list ...interface{}) error {
	return errNotImplemented
}
func (*sqlXWrap) Update(list ...interface{}) (int64, error) {
	return 0, errNotImplemented

}
func (*sqlXWrap) Delete(list ...interface{}) (int64, error) {
	return 0, errNotImplemented

}
func (*sqlXWrap) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, errNotImplemented

}
func (*sqlXWrap) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	err := sqlxDB.Select(i, query, args...)
	return nil, err
}
func (*sqlXWrap) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return sqlxDB.Get(holder, query, args...)
}

var errNotImplemented error = errors.New("Not implemented")
