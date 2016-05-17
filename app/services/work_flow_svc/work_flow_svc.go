package work_flow_svc

import (
	"github.com/aws/aws-sdk-go/service/marketplacecommerceanalytics"
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type IWorkFlowService interface {
	Start()
}

type irepo interface {
	GetFailedToReachConsensusQuestions() ([]Question, error)
}

type Question struct {
	models.Question
	MicroTask models.MicroTask
}

type workFlowService struct {
	r irepo
}

func (s *workFlowService) Start() {
	questions, err := s.r.GetFailedToReachConsensusQuestions()
	if err != nil {
		return
	}

}
