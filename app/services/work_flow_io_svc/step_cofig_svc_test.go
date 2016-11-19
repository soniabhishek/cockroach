package work_flow_io_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stepRepoMock struct {
	step_repo.IStepRepo
	StepToReturn models.Step
}

func (s *stepRepoMock) GetById(stepId uuid.UUID) (models.Step, error) {
	return s.StepToReturn, nil
}

func TestStepConfigSvc_GetCrowdsourcingStepConfig(t *testing.T) {

	stepRepo := &stepRepoMock{IStepRepo: step_repo.Mock()}

	stepConfigSvc := stepConfigSvc{
		stepRepo: stepRepo,
	}

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{}}
	cc, err := stepConfigSvc.GetCrowdsourcingStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{microTaskId: "321", answerKey: "abcd"}}
	cc, err = stepConfigSvc.GetCrowdsourcingStepConfig(uuid.NewV4())
	assert.Error(t, err)

	mId := uuid.NewV4()
	stepRepo.StepToReturn = models.Step{Config: models.JsonF{microTaskId: mId.String(), answerKey: "abcd"}}
	cc, err = stepConfigSvc.GetCrowdsourcingStepConfig(uuid.NewV4())
	assert.NoError(t, err)
	assert.EqualValues(t, mId, cc.MicroTaskId)
	assert.EqualValues(t, "abcd", cc.AnswerKey)
}

func TestStepConfigSvc_GetTransformationStepConfig(t *testing.T) {

	stepRepo := &stepRepoMock{IStepRepo: step_repo.Mock()}

	stepConfigSvc := stepConfigSvc{
		stepRepo: stepRepo,
	}

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{}}
	tc, err := stepConfigSvc.GetTransformationStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{templateId: "321"}}
	tc, err = stepConfigSvc.GetTransformationStepConfig(uuid.NewV4())
	assert.NoError(t, err)
	assert.EqualValues(t, "321", tc.TemplateId)

}

func TestStepConfigSvc_GetBifurcationStepConfig(t *testing.T) {

	stepRepo := &stepRepoMock{IStepRepo: step_repo.New()}

	stepConfigSvc := stepConfigSvc{
		stepRepo: stepRepo,
	}

	stepRepo.StepToReturn = models.Step{}
	bc, err := stepConfigSvc.GetBifurcationStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{}}
	bc, err = stepConfigSvc.GetBifurcationStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{multiplication: "321"}}
	bc, err = stepConfigSvc.GetBifurcationStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{multiplication: 2}}
	bc, err = stepConfigSvc.GetBifurcationStepConfig(uuid.NewV4())
	assert.NoError(t, err)
	assert.EqualValues(t, 2, bc.Multiplication)

}

func TestStepConfigSvc_GetAlgorithmStepConfig(t *testing.T) {

	stepRepo := &stepRepoMock{IStepRepo: step_repo.New()}

	stepConfigSvc := stepConfigSvc{
		stepRepo: stepRepo,
	}

	stepRepo.StepToReturn = models.Step{}
	bc, err := stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{answerKey: 121, textFieldKey: "text_field"}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{answerKey: "", textFieldKey: "text_field"}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{textFieldKey: "text_field"}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.Error(t, err)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{answerKey: "answer_field", textFieldKey: "text_field", timeDelayStart: 3, timeDelayStop: 4}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.NoError(t, err)
	assert.EqualValues(t, "answer_field", bc.AnswerKey)
	assert.EqualValues(t, "text_field", bc.TextFieldKey)
	assert.EqualValues(t, 3, bc.TimeDelayStart)
	assert.EqualValues(t, 4, bc.TimeDelayStop)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{answerKey: "answer_field", textFieldKey: "text_field"}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.NoError(t, err)
	assert.EqualValues(t, 0, bc.TimeDelayStart)

	stepRepo.StepToReturn = models.Step{Config: models.JsonF{answerKey: "answer_field", textFieldKey: "text_field", timeDelayStart: "23"}}
	bc, err = stepConfigSvc.GetAlgorithmStepConfig(uuid.NewV4())
	assert.NoError(t, err)
	assert.EqualValues(t, 0, bc.TimeDelayStart)
	assert.EqualValues(t, 0, bc.TimeDelayStop)
}
