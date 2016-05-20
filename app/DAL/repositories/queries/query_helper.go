package queries

import sq "github.com/Masterminds/squirrel"

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
