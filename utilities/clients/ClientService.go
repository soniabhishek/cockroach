package main

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type client struct {
	userName       string
	password       []byte
	clientSecretId uuid.UUID
	projectId      uuid.UUID
	projectName    string
	projectDesc    string
	url            string
	header         models.JsonFake
	steps          []int
}

func createClient(cl client) (status bool, err error) {

	return
}
