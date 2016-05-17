package macro_task_repo

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/data_access_svc/clients"
)

type IMacroTaskRepo interface {
	Get(uuid.UUID) (models.MacroTask, error)
}

func New() IMacroTaskRepo {
	return &macroTaskRepo{
		pg:  clients.GetPostgresClient(),
		mgo: clients.GetMongoClient(),
	}
}
