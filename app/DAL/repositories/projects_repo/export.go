package projects_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IProjectsRepo interface {
	Get(uuid.UUID) (models.Project, error)
}

func New() IProjectsRepo {
	return &projectsRepo{
		pg:  postgres.GetPostgresClient(),
		mgo: clients.GetMongoClient(),
	}
}
