package repositories

import "database/sql"

type IDatabase interface {
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
	SelectById(holder interface{}, id string) error

	//SelectInt(query string, args ...interface{}) (int64, error)
	//SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)
	//SelectFloat(query string, args ...interface{}) (float64, error)
	//SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error)
	//SelectStr(query string, args ...interface{}) (string, error)
	//SelectNullStr(query string, args ...interface{}) (sql.NullString, error)
}
