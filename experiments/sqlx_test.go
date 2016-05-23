package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"reflect"
	"testing"
)

func TestSqlXWrap_SelectCustom(t *testing.T) {
	mId := uuid.FromStringOrNil("446ac588-6f12-41ba-b93b-ce5c077fdc67")
	maId := uuid.FromStringOrNil("ce386fe6-409d-45a0-8b9f-eb9877de3923")

	g := postgres.GetGorpClient()

	r, err := g.Db.Query(`select ma.*,m.* from micro_tasks m inner join
	 macro_tasks ma on m.macro_task_id = ma.id
	 where m.id = '` + mId.String() + `'`)
	assert.NoError(t, err)

	type MicroTaskExt struct {
		models.MacroTask
		models.MicroTask
	}

	m := MicroTaskExt{}

	var dest interface{} = &m

	for r.Next() {
		v := reflect.ValueOf(dest)

		if v.Kind() != reflect.Ptr {
			err = fmt.Errorf("must pass a pointer, not a value, to StructScan destination")
			assert.NoError(t, err)
		}

		v = reflect.Indirect(v)
		mrT := reflect.TypeOf(&m).Elem()

		columns, err := r.Columns()
		assert.NoError(t, err)

		rvalues := make([]interface{}, len(columns))
		start := 0
		end := 0

		for i := 0; i < v.NumField(); i++ {

			inner := mrT.Field(i).Type
			innerV := v.Field(i)
			fieldTags := getDbFieldTagMap(innerV, inner)

			start = end
			end += len(fieldTags)
			innerCols := columns[start:end]

			for tag, field := range fieldTags {
				ind := getColPosition(innerCols, tag)
				rvalues[start+ind] = field.Addr().Interface()
			}

		}

		err = r.Scan(rvalues...)

		assert.NoError(t, err)
	}

	assert.EqualValues(t, mId.String(), m.MicroTask.ID.String())
	assert.EqualValues(t, maId.String(), m.MacroTask.ID.String())
}

func getColPosition(cols []string, name string) int {
	for i, v := range cols {
		if name == v {
			return i
		}
	}
	return -1
}

func getDbFieldTagMap(v reflect.Value, t reflect.Type) map[string]reflect.Value {

	fieldTag := make(map[string]reflect.Value, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		typ := t.Field(i)
		val := v.Field(i)
		if dbTag := typ.Tag.Get("db"); dbTag != "" {
			fieldTag[dbTag] = val
		}
	}
	return fieldTag
}
