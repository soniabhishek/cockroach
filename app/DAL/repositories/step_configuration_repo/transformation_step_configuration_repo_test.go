package step_configuration_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformationStepConfigRepo_GetByStepId(t *testing.T) {
	t.SkipNow()
	tStepConfigRepo := transformationStepConfigRepo{Db: postgres.GetPostgresClient()}
	result, err := tStepConfigRepo.GetByStepId(uuid.FromStringOrNil("05da5599-221e-4973-854a-9fe39b12339e"))
	assert.NoError(t, err)
	assert.EqualValues(t, uuid.FromStringOrNil("05da5599-221e-4973-854a-9fe39b12339e"), result.StepId)
}
