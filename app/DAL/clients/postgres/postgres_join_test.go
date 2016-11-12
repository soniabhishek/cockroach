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
		models.User
		models.Client
	}

	var macroUser []MacroTaskWithCreator

	err := pg.SelectOneJoin(&macroUser, `select u.*, c.* from users u, clients c order by u.username`)
	fmt.Println("FINAL OUTPUT", macroUser)
	assert.NoError(t, err)
	//assert.EqualValues(t, macroUser.User.ID.String(), macroUser.CreatorId.String())
}
