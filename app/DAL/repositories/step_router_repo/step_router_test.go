package step_router_repo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func TestInMemStepRouteRepo_GetRoutesWithLogicByStepId(t *testing.T) {
	pgClient := postgres.GetPostgresClient()

	sr := stepRouteRepo{
		Db: pgClient,
	}

	rs, err := sr.GetRoutesWithLogicByStepId(uuid.FromStringOrNil("d39184f4-aa64-4f13-b7b5-88a72bd6d68c"))
	assert.NoError(t, err)
	fmt.Println(rs)
}
