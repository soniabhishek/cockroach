package repositories

import (
	"database/sql"
	"fmt"
)

type IDatabase interface {
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
	SelectById(holder interface{}, id fmt.Stringer) error
	SelectOneJoin(holder interface{}, query string, args ...interface{}) error
}
