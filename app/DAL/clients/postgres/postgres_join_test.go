package postgres_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"testing"
)

func TestPostgres_db_SelectOneJoin(t *testing.T) {
	pg := postgres.GetPostgresClient()

	type MacroTaskWithCreator struct {
		models.MacroTask
		models.User
	}

	var macroUser MacroTaskWithCreator

	err := pg.SelectOneJoin(&macroUser, `select m.*,u.* from macro_tasks m inner join
	users u on u.id = m.creator_id limit 1`)
	assert.NoError(t, err)
	assert.EqualValues(t, macroUser.User.ID.String(), macroUser.CreatorId.String())
}
