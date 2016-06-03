package projects_repo

import (
	"database/sql"

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

func (r *projectsRepo) Get(id uuid.UUID) (m models.Project, err error) {
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
