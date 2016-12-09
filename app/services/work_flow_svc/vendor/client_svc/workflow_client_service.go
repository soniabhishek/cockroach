package client_svc

import (
	"database/sql"
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/pkg/errors"
)

type workFlowClientService struct {
	clientsRepo clients_repo.IClientsRepo
	userRepo    user_repo.IUserRepo
}

func (wc *workFlowClientService) CreateClient(client models.Client) (response models.Client, err error) {
	if client.Name == "" {
		err = errors.New("Invalid Client Name")
		return
	}
	user := models.User{}
	user.Username = client.Name
	user.FirstName = sql.NullString{client.Name, true}

	err = wc.userRepo.Add(&user)
	if err != nil {
		return
	}
	client.UserId = user.ID
	err = wc.clientsRepo.Add(&client)
	if err != nil {
		return
	}
	client.ClientSecretUuid = uuid.Nil
	return client, nil
}

func (wc *workFlowClientService) GetClient(uuid.UUID) (client models.Client, err error) {
	return
}

func (wc *workFlowClientService) FetchAllClient() ([]models.Client, error) {
	return wc.clientsRepo.GetAllClients()
}
