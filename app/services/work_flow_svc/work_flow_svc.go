package work_flow_svc

import (
	"errors"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"strings"
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
	textFieldKey   = "text_field_key"
	timeDelayStart = "time_delay_start"
	timeDelayStop  = "time_delay_stop"
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
	//TODO: Change the trimspace logic to where the data is getting inserted in db
	microTaskUuid, err := uuid.FromString(strings.TrimSpace(microTaskIdString))
	if err != nil {
		err = ErrConfigNotFound
		return
	}
	tc.MicroTaskId = microTaskUuid
	tc.AnswerKey = strings.TrimSpace(answerKeyString)
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
	tc.TemplateId = strings.TrimSpace(templateIdString)
	return
}

func (s *stepConfigSvc) GetValidationStepConfig(stepId uuid.UUID) (tc models.ValidationConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	templateID, ok := step.Config[templateId]
	answerKey, ok2 := step.Config[answerKey]
	if !ok || !ok2 {
		err = ErrConfigNotFound
		return
	}
	templateIdString, ok := templateID.(string)
	answerKeyString, ok2 := answerKey.(string)
	if !ok || !ok2 {
		err = ErrConfigNotFound
		return
	}
	tc.TemplateId = strings.TrimSpace(templateIdString)
	tc.TemplateId = strings.TrimSpace(answerKeyString)
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
func (s *stepConfigSvc) GetAlgorithmStepConfig(stepId uuid.UUID) (ac models.AlgorithmConfig, err error) {
	step, err := s.stepRepo.GetById(stepId)
	if err != nil {
		return
	}

	answerFieldKey, ok := step.Config[answerKey]
	textFieldKey, ok2 := step.Config[textFieldKey]
	if !ok || !ok2 {
		err = ErrConfigNotFound
		return
	}
	answerFieldKeyString, ok := answerFieldKey.(string)
	textFieldKeyString, ok2 := textFieldKey.(string)
	if !ok || !ok2 || answerFieldKey == "" || textFieldKey == "" {
		err = ErrConfigNotFound
		return
	}

	timeDelayStart, ok := step.Config[timeDelayStart]
	if ok {
		timeDelayStartFloat, ok1 := timeDelayStart.(float64)

		if ok1 {
			ac.TimeDelayStart = timeDelayStartFloat
		} else {
			ac.TimeDelayStart = 0
		}
	} else {
		ac.TimeDelayStart = 0
	}

	timeDelayStop, ok := step.Config[timeDelayStop]
	if ok {
		timeDelayStopFloat, ok1 := timeDelayStop.(float64)
		if ok1 {
			ac.TimeDelayStop = timeDelayStopFloat
		} else {
			ac.TimeDelayStop = 0
		}
	} else {
		ac.TimeDelayStop = 0
	}

	ac.AnswerKey = strings.TrimSpace(answerFieldKeyString)
	ac.TextFieldKey = strings.TrimSpace(textFieldKeyString)

	return
}
