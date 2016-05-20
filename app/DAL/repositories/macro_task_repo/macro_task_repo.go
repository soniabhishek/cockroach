package macro_task_repo

import (
	"database/sql"

	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/DAL/repositories/queries"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gopkg.in/mgo.v2"
)

const macroTaskTable = "macro_tasks"

type macroTaskRepo struct {
	pg  repositories.IDatabase
	mgo *mgo.Database
}

func (r *macroTaskRepo) getFromPG(id uuid.UUID) (m models.MacroTask, err error) {
	err = r.pg.SelectOne(&m, queries.SelectById(macroTaskTable), id)
	err = transformErr(err)
	return
}

func (r *macroTaskRepo) getFromMgo(id uuid.UUID) (m models.MacroTask, err error) {
	err = r.mgo.C("macro_tasks").FindId(id).One(&m)
	err = transformErr(err)
	return
}

func (r *macroTaskRepo) Get(id uuid.UUID) (m models.MacroTask, err error) {
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

func (r *macroTaskRepo) saveMgo(m models.MacroTask) error {
	return r.mgo.C("macro_tasks").Insert(&m)
}

func transformErr(err error) error {

	switch err {
	case sql.ErrNoRows:
		fallthrough
	case mgo.ErrNotFound:
		err = ErrMacroTaskNotFound
	}
	return err
}
