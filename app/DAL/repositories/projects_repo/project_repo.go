package projects_repo

import (
	"database/sql"

	"errors"
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/DAL/repositories/queries"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"gopkg.in/mgo.v2"
	"time"
)

const projectTable = "projects"

type projectsRepo struct {
	pg  repositories.IDatabase
	mgo *mgo.Database
}

var _ IProjectsRepo = (*projectsRepo)(nil)

func (r *projectsRepo) getFromPG(id uuid.UUID) (m models.Project, err error) {
	err = r.pg.SelectOne(&m, queries.SelectById(projectTable), id)
	err = transformErr(err)
	return
}

func (r *projectsRepo) getFromMgo(id uuid.UUID) (m models.Project, err error) {
	err = r.mgo.C(projectTable).FindId(id).One(&m)
	err = transformErr(err)
	return
}

func (r *projectsRepo) GetById(id uuid.UUID) (m models.Project, err error) {
	m, err = r.getFromMgo(id)
	if err == nil {
		return
	}
	m, err = r.getFromPG(id)
	if err != nil {
		return
	}
	go r.saveMgo(m)
	return
}

func (r *projectsRepo) GetByClientId(id uuid.UUID) (m []models.Project, err error) {
	_, err = r.pg.Select(&m, `select * from projects where client_id = $1`, id.String())
	return
}

func (r *projectsRepo) saveMgo(m models.Project) error {
	return r.mgo.C(projectTable).Insert(&m)
}

func transformErr(err error) error {

	switch err {
	case sql.ErrNoRows:
		fallthrough
	case mgo.ErrNotFound:
		err = ErrProjectNotFound
	}
	return err
}

func (i *projectsRepo) Add(p *models.Project) error {
	p.ID = uuid.NewV4()
	p.UpdatedAt = pq.NullTime{time.Now(), true}
	p.CreatedAt = p.UpdatedAt
	p.StartedAt = p.UpdatedAt
	p.EndedAt = pq.NullTime{time.Time{}, false}
	return i.pg.Insert(p)
}
func (i *projectsRepo) Update(p models.Project) error {
	_, err := i.pg.Update(&p)
	return err
}
func (i *projectsRepo) Delete(id uuid.UUID) error {
	query := fmt.Sprintf(`delete from projects where id='%v'::uuid`, id)
	res, err := i.pg.Exec(query)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		err = errors.New("Could not delete Client with ID [" + id.String() + "]")
	}
	return err
}
func (i *projectsRepo) IfIdExist(id uuid.UUID) (ifExist bool, err error) {
	err = i.pg.SelectOne(&ifExist, `select exists(select 1 from projects where id=$1)`, id)
	if err != nil {
		return
	}
	return
}
