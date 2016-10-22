package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/DAL/repositories/queries"
	"gopkg.in/gorp.v1"
)

//This wrapper postgres_db is to have one place where we can write all the modifications on errors thrown by gorp & pq

type postgres_db struct {
	gorpDbMap *gorp.DbMap
}

var _ repositories.IDatabase = &postgres_db{}

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

func (pg *postgres_db) SelectById(holder interface{}, id fmt.Stringer) error {

	ptrv := reflect.ValueOf(holder)
	if ptrv.Kind() != reflect.Ptr {
		e := fmt.Sprintf("passed non-pointer: %v (kind=%v)", holder,
			ptrv.Kind())
		return errors.New(e)
	}
	elem := ptrv.Elem()
	etype := reflect.TypeOf(elem.Interface())

	table, err := pg.gorpDbMap.TableFor(etype, false)
	if err != nil {
		return err
	}
	return pg.gorpDbMap.SelectOne(holder, queries.SelectById(table.TableName), id.String())
}

func GetPostgresClient() *postgres_db {
	return &postgres_db{gorpDbMap}
}
