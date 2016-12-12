package postgres

import (
	"database/sql"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/plog"
	"gopkg.in/gorp.v1"
)

//This wrapper postgres_db is to have one place where we can write all the modifications on errors thrown by gorp & pq

type transactionalPostgres struct {
	gorpDbMap *gorp.DbMap
	trans     *gorp.Transaction
}

var _ repositories.IDatabase = &transactionalPostgres{}

func (tp *transactionalPostgres) Insert(list ...interface{}) (err error) {
	return tp.trans.Insert(list...)
}
func (tp *transactionalPostgres) Update(list ...interface{}) (d int64, err error) {
	return tp.trans.Update(list...)
}
func (pg *transactionalPostgres) Delete(list ...interface{}) (int64, error) {
	return pg.trans.Delete(list...)
}
func (pg *transactionalPostgres) Exec(query string, args ...interface{}) (sql.Result, error) {
	return pg.trans.Exec(query, args...)
}
func (pg *transactionalPostgres) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return pg.trans.Select(i, query, args...)
}
func (pg *transactionalPostgres) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return pg.trans.SelectOne(holder, query, args...)
}

func (pg *transactionalPostgres) SelectById(holder interface{}, id fmt.Stringer) error {
	return selectById(pg.gorpDbMap.TableFor, pg.trans.SelectOne, holder, id)
}

//TODO remove usage of Db.Query here
func (pg *transactionalPostgres) SelectOneJoin(holder interface{}, query string, args ...interface{}) error {

	return selectOneJoin(pg.gorpDbMap.Db.Query, holder, query, args...)
}

func (pg *transactionalPostgres) SelectJoin(holder interface{}, query string, args ...interface{}) error {

	return selectJoin(pg.gorpDbMap.Db.Query, holder, query, args...)
}

func (pg *transactionalPostgres) Commit() {
	if err := pg.trans.Commit(); err != nil {
		plog.Error("Postgres client", err, "Error occured while Commit transaction")
		panic(err)
	}
}

func (pg *transactionalPostgres) Rollback() {
	if err := pg.trans.Rollback(); err != nil {
		plog.Error("Postgres client", err, "Error occured in Rollback transaction")
		panic(err)
	}
}
func GetTransactionClient() *transactionalPostgres {

	tx, err := gorpDbMap.Begin()
	if err != nil {
		plog.Error("Postgres client", err, "Error occured in creating transaction")
		panic(err)
	}

	return &transactionalPostgres{gorpDbMap, tx}
}
