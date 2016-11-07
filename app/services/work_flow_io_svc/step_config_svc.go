package work_flow_io_svc

import (
	"errors"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
)

type stepConfigSvc struct {
	stepRepo step_repo.IStepRepo
}

var ErrConfigNotFound = errors.New("Configuration not found")
var ErrConfigMalformed = errors.New("Configuration Malformed")

var _ IStepConfigSvc = &stepConfigSvc{}

const (
	templateId     = "template_id"
	multiplication = "multiplication"
	microTaskId    = "micro_task_id"
	answerKey      = "answer_key"
	imageFieldKey  = "image_field_key"
)

func (s *stepConfigSvc) GetCrowdsourcingStepConfig(stepId uuid.UUID) (tc models.CrowdsourcingConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	microTaskId, ok := step.Config[microTaskId]
	answerKey, ok1 := step.Config[answerKey]
	if !ok || !ok1 {
		err = ErrConfigNotFound
		return
	}
	microTaskIdString, ok := microTaskId.(string)
	answerKeyString, ok1 := answerKey.(string)
	if !ok || !ok1 {
		err = ErrConfigNotFound
		return
	}

	microTaskUuid, err := uuid.FromString(microTaskIdString)
	if err != nil {
		err = ErrConfigNotFound
		return
	}
	tc.MicroTaskId = microTaskUuid
	tc.AnswerKey = answerKeyString
	return
}

func (s *stepConfigSvc) GetTransformationStepConfig(stepId uuid.UUID) (tc models.TransformationConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	templateID, ok := step.Config[templateId]
	if !ok {
		err = ErrConfigNotFound
		return
	}
	templateIdString, ok := templateID.(string)
	if !ok {
		err = ErrConfigNotFound
		return
	}
	tc.TemplateId = templateIdString
	return
}

func (s *stepConfigSvc) GetBifurcationStepConfig(stepId uuid.UUID) (bc models.BifurcationConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	err = step.Config.CastTo(&bc)
	if err != nil {
		return
	}

	if bc.Multiplication < 1 {
		err = ErrConfigNotFound
		plog.Error("StepConfigSvc", ErrConfigMalformed, "stepId "+stepId.String())
		return
	}

	return
}

func (s *stepConfigSvc) GetUnificationStepConfig(stepId uuid.UUID) (uc models.UnificationConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	err = step.Config.CastTo(&uc)
	if err != nil {
		return
	}

	if uc.Multiplication < 1 {
		err = ErrConfigNotFound
		plog.Error("StepConfigSvc", ErrConfigMalformed, "stepId "+stepId.String())
		return
	}

	return
}

func (s *stepConfigSvc) GetStartStepConfig(stepId uuid.UUID) (uc models.StartStepConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	imageFieldKey, ok := step.Config[imageFieldKey]
	if !ok {
		err = ErrConfigNotFound
		return
	}

	imageFieldKeyString := imageFieldKey.(string)

	uc.ImageFieldKey = imageFieldKeyString

	return
}
