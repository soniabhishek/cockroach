package postgres_test

import (
	"testing"

	"fmt"
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_db_SelectOneJoin(t *testing.T) {
	pg := postgres.GetPostgresClient()

	type MacroTaskWithCreator struct {
		models.Client
		models.User
	}

	var macroUser MacroTaskWithCreator

	err := pg.SelectOneJoin(&macroUser, `select c.*,u.* from clients c inner join
	users u on u.id = c.user_id limit 1`)
	fmt.Println("SOLUTION", macroUser)
	assert.NoError(t, err)
	//assert.EqualValues(t, macroUser.User.ID.String(), macroUser.CreatorId.String())
}
