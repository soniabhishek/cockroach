package client_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkFlowClientService interface {
	CreateClient(models.Client) (models.Client, error)
	GetClient(uuid.UUID) (models.Client, error)
	FetchAllClient() ([]models.Client, error)
}

func New() IWorkFlowClientService {
	return &workFlowClientService{
		clientsRepo: clients_repo.New(),
		userRepo:    user_repo.New(),
	}
}
