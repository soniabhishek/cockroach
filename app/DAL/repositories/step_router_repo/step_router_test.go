package step_router_repo

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
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
