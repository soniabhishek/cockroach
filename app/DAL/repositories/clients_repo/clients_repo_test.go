package clients_repo

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClientsRepo_GetByProjectId(t *testing.T) {

	t.SkipNow()

	pgClient := postgres.GetPostgresClient()
	cl := clientsRepo{
		Db: pgClient,
	}

	// make it pick project from db
	client, err := cl.GetByProjectId(uuid.FromStringOrNil("4aa64098-c74a-4efc-adf9-061ce3c70bf4"))
	assert.Error(t, err)
	fmt.Println(client)

	client, err = cl.GetByProjectId(uuid.FromStringOrNil("416f88d2-4172-43c7-abab-555f56bf656b"))
	assert.NoError(t, err)
	fmt.Println(client)
}
