package step_configuration_repo

import (
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type transformationStepConfigRepo struct {
	Db repositories.IDatabase
}

var _ ITransformationStepConfigurationRepo = &transformationStepConfigRepo{}

const transformationStepConfigTable string = "transformation_step_configuration"

func (t *transformationStepConfigRepo) GetByStepId(stepId uuid.UUID) (models.TransformationStepConfiguration, error) {

	var transformationDetails models.TransformationStepConfiguration

	query := `SELECT * FROM ` + transformationStepConfigTable + ` WHERE step_id = $1`

	err := t.Db.SelectOne(&transformationDetails, query, stepId)

	if err != nil {
		return transformationDetails, err
	}
	return transformationDetails, nil
}
