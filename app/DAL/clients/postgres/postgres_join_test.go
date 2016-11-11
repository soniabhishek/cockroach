package postgres_test

import (
	"testing"

	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/stretchr/testify/assert"
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
