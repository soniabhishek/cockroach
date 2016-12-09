package clients_repo

import (
	"errors"
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"time"
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

func (c *clientsRepo) Add(cl *models.Client) error {
	cl.ID = uuid.NewV4()
	cl.ClientSecretUuid = uuid.NewV4()
	cl.CreatedAt = pq.NullTime{time.Now(), true}
	cl.UpdatedAt = cl.CreatedAt
	return c.Db.Insert(cl)
}

func (c *clientsRepo) Update(cl models.Client) error {
	_, err := c.Db.Update(&cl)
	return err
}

func (c *clientsRepo) Delete(id uuid.UUID) error {
	query := fmt.Sprintf(`delete from clients where id='%v'::uuid`, id)
	res, err := c.Db.Exec(query)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		err = errors.New("Could not delete Client with ID [" + id.String() + "]")
	}
	return err
}

func (c *clientsRepo) GetAllClients() (clients []models.Client, err error) {
	_, err = c.Db.Select(&clients, `select id, name from clients`)
	return
}
func (c *clientsRepo) IfIdExist(id uuid.UUID) (ifExist bool, err error) {
	err = c.Db.SelectOne(&ifExist, `select exists(select 1 from clients where id=$1)`, id)
	return
}
