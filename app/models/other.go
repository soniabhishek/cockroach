package models

import (
	"database/sql"

	"gitlab.com/playment-main/support/app/models/uuid"
	"gopkg.in/gorp.v1"
)

type ImageContainer struct {
	Id  string
	Url string
}

type ImageDictionary struct {
	ID             uuid.UUID
	Label          string
	OriginalUrl    string
	IsCached       bool
	CdnUrl         string
	ResizedCdnUrl  string
	CreatedAt      gorp.NullTime
	BatchProcessId uuid.UUID
}

type ImageDictionary1 struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	RealUrl   string         `db:"real_url" json:"real_url"`
	CloudUrl  string         `db:"cloud_url" json:"cloud_url"`
	Extra     sql.NullString `db:"extra" json:"extra"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type BatchProcess struct {
	ID          uuid.UUID
	Name        string
	Done        bool
	Completion  float64
	Type        int8
	CreatedAt   gorp.NullTime
	MacroTaskId uuid.UUID
}
