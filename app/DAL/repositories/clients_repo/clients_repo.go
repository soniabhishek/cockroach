package clients_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type clientsRepo struct {
	Db repositories.IDatabase
}

func (c *clientsRepo) GetByProjectId(projectId uuid.UUID) (models.Client, error) {
	var client models.Client

	err := c.Db.SelectOne(&client, `select c.* from clients c
	inner join projects p on p.client_id  = c.id
	where p.id = $1
	`, projectId)
	return client, err
}
