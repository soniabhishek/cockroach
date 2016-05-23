package macro_task_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IMacroTaskRepo interface {
	Get(uuid.UUID) (models.MacroTask, error)
}

func New() IMacroTaskRepo {
	return &macroTaskRepo{
		pg:  postgres.GetPostgresClient(),
		mgo: clients.GetMongoClient(),
	}
}
