package postgres

import (
	"errors"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/repositories/queries"
	"gopkg.in/gorp.v1"
	"reflect"
)

type tableForFunc func(t reflect.Type, checkPK bool) (*gorp.TableMap, error)
type selectOneFunc func(holder interface{}, query string, args ...interface{}) error

func selectById(tableFor tableForFunc, selectOne selectOneFunc, holder interface{}, id fmt.Stringer) error {

	ptrv := reflect.ValueOf(holder)
	if ptrv.Kind() != reflect.Ptr {
		e := fmt.Sprintf("passed non-pointer: %v (kind=%v)", holder,
			ptrv.Kind())
		return errors.New(e)
	}
	elem := ptrv.Elem()
	etype := reflect.TypeOf(elem.Interface())

	table, err := tableFor(etype, false)
	if err != nil {
		return err
	}
	return selectOne(holder, queries.SelectById(table.TableName), id.String())
}
