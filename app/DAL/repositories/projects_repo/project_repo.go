package projects_repo

import (
	"database/sql"

	"errors"
	"fmt"

	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/DAL/repositories/queries"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gopkg.in/mgo.v2"
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

func (i *projectsRepo) Add(p models.Project) error {
	return i.pg.Insert(&p)
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
