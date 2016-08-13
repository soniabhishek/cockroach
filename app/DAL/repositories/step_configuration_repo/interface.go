package step_configuration_repo

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type ITransformationStepConfigurationRepo interface {
	GetByStepId(stepId uuid.UUID) (models.TransformationStepConfiguration, error)
}
