package projects_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
)

func New() IProjectsRepo {
	return &projectsRepo{
		pg:  postgres.GetPostgresClient(),
		mgo: clients.GetMongoClient(),
	}
}
