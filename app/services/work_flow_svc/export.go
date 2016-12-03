package work_flow_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IStepConfigSvc interface {
	GetCrowdsourcingStepConfig(stepId uuid.UUID) (models.CrowdsourcingConfig, error)
	GetTransformationStepConfig(stepId uuid.UUID) (models.TransformationConfig, error)
	GetBifurcationStepConfig(stepId uuid.UUID) (models.BifurcationConfig, error)
	GetUnificationStepConfig(stepId uuid.UUID) (models.UnificationConfig, error)
	GetAlgorithmStepConfig(stepId uuid.UUID) (models.AlgorithmConfig, error)
}

func NewStepConfigService() IStepConfigSvc {
	return &stepConfigSvc{
		stepRepo: step_repo.New(),
	}
}
