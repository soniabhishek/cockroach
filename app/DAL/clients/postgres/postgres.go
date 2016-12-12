package postgres

import (
	"database/sql"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/repositories"
	"gopkg.in/gorp.v1"
)

//This wrapper postgres_db is to have one place where we can write all the modifications on errors thrown by gorp & pq

type postgresDb struct {
	gorpDbMap *gorp.DbMap
}

var _ repositories.IDatabase = &postgresDb{}

func (pg *postgresDb) Insert(list ...interface{}) (err error) {
	return pg.gorpDbMap.Insert(list...)
}
func (pg *postgresDb) Update(list ...interface{}) (d int64, err error) {
	return pg.gorpDbMap.Update(list...)
}
func (pg *postgresDb) Delete(list ...interface{}) (int64, error) {
	return pg.gorpDbMap.Delete(list...)
}
func (pg *postgresDb) Exec(query string, args ...interface{}) (sql.Result, error) {
	return pg.gorpDbMap.Exec(query, args...)
}
func (pg *postgresDb) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return pg.gorpDbMap.Select(i, query, args...)
}
func (pg *postgresDb) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return pg.gorpDbMap.SelectOne(holder, query, args...)
}

func (pg *postgresDb) SelectById(holder interface{}, id fmt.Stringer) error {

	return selectById(pg.gorpDbMap.TableFor, pg.gorpDbMap.SelectOne, holder, id)
}

func (pg *postgresDb) SelectOneJoin(holder interface{}, query string, args ...interface{}) error {

	return selectOneJoin(pg.gorpDbMap.Db.Query, holder, query, args...)
}

func (pg *postgresDb) SelectJoin(holder interface{}, query string, args ...interface{}) error {

	return selectJoin(pg.gorpDbMap.Db.Query, holder, query, args...)
}

func GetPostgresClient() *postgresDb {
	return &postgresDb{gorpDbMap}
}
