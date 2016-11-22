package models

import "github.com/crowdflux/angel/app/models/uuid"

type ClientModel struct {
	Id       uuid.UUID `db:"id" json:"id"`
	UserName string    `db:"username" json:"name"`
}
