package projects_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
)

func New() IProjectsRepo {
	return &projectsRepo{
		pg:  postgres.GetPostgresClient(),
		mgo: clients.GetMongoClient(),
	}
}
