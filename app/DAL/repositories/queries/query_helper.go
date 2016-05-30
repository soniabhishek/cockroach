package queries

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func SelectById(tableName string) (sql string) {

	sql, _, err := sq.Select("*").
		From(tableName).
		Where("id = $1").
		ToSql()

	if err != nil {
		panic(err)
	}
	return
}

func SelectByName(tableName string) (sql string) {

	sql, _, err := sq.Select("*").
		From(tableName).
		Where("name = $1").
		ToSql()

	if err != nil {
		panic(err)
	}
	return
}

// Returns a new UUID if input is Nil otherwise the same
func EnsureId(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.NewV4()
	}
	return id
}
