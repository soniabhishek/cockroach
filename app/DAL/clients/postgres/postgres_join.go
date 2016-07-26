package postgres

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/crowdflux/angel/utilities"
)

const dbTag string = "db"

// This method allows you to write joins & get results from
// multiple tables in one query
//
// The holder should embed all the structs for which the query
// is returning the data IN ORDER
//
// e.g. pointer to a struct MicroTaskExt can be passed as holder
// type MicroTaskExt struct {
//	models.MacroTask
// 	models.MicroTask
// }
// and input query must give the results in form of macro_tasks.*, micro_task.*
// and NO EXTRA crap (see postgres_join_test.go for more)
//
// It will fill the given struct with values
// and can be taken out for other use
//
// TODO Currently this doesn't support multi level embedded structs. Add it in future
func (pg *postgres_db) SelectOneJoin(holder interface{}, query string, args ...interface{}) error {

	//Db.query --> native sql drive's query method
	r, err := pg.gorpDbMap.Db.Query(query, args...)
	if err != nil {
		return err
	}

	//check if any row is present
	if r.Next() {

		// reflect value of holder
		v := reflect.ValueOf(holder)

		// value should be a pointer
		if v.Kind() != reflect.Ptr {
			err = fmt.Errorf("must pass a pointer, not a value, to StructScan destination")
			return err
		}

		// Get the underlying val of pointer
		v = reflect.Indirect(v)

		// get column names in order according to sql query
		columns, err := r.Columns()
		if err != nil {
			return err
		}

		rvalues := mapColumnsToInterfaces(v, columns)

		// Scan method takes interfaces & puts values in it according to the order
		return r.Scan(rvalues...)

	} else {
		return errors.New("No matching rows found")
	}
}

func getDbTagFieldMap(v reflect.Value, t reflect.Type) map[string]reflect.Value {

	fieldCount := v.NumField()

	fieldTag := map[string]reflect.Value{}

	for i := 0; i < fieldCount; i++ {
		typ := t.Field(i)
		val := v.Field(i)

		if dbTag := typ.Tag.Get(dbTag); dbTag != "" && dbTag != "-" {
			fieldTag[dbTag] = val
		}
	}
	return fieldTag
}

func mapColumnsToInterfaces(inVal reflect.Value, columns []string) []interface{} {

	// Get the underlying val's type
	inTyp := inVal.Type()

	// initial the array to hold references of struct fields mapped to the columns
	rvalues := make([]interface{}, len(columns))
	start := 0
	end := 0

	// Loop on the embedded structs
	for i := 0; i < inVal.NumField(); i++ {

		// Get ith embedded type & value
		embeddedT := inTyp.Field(i).Type
		embeddedV := inVal.Field(i)

		// Get dbTag-Field Map of the struct
		tagFieldMap := getDbTagFieldMap(embeddedV, embeddedT)

		// start from the last end
		start = end

		// get embedded struct's number of fields and put the end there
		end += len(tagFieldMap)

		// filter all columns by start & end
		innerCols := columns[start:end]

		// Loop over the tag field map & put the reference in rvalues array
		for tag, field := range tagFieldMap {
			if ind := utilities.FindIndexStringArr(innerCols, tag); ind >= 0 {
				rvalues[start+ind] = field.Addr().Interface()
			} else {
				// handle case when db tag is not in the column
			}
		}
	}
	return rvalues
}
